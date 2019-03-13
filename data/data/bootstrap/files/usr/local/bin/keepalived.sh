#!/usr/bin/env bash
set -e

function get_iface_in_vip_subnet() {
    local vip
    local net_cidr
    local iface_cidrs

    vip="$1"

    iface_cidrs=$(ip addr show | grep -v "scope host" | grep -Po 'inet \K[\d.]+/[\d.]+' | xargs)
    net_cidr=$(/usr/bin/env python - "$vip" "$iface_cidrs" << EOF
import sys
import socket
import struct

vip = sys.argv[1]
iface_cidrs = sys.argv[2].split()
vip_int = struct.unpack("!I", socket.inet_aton(vip))[0]

for iface_cidr in iface_cidrs:
    ip, prefix = iface_cidr.split('/')
    ip_int = struct.unpack("!I", socket.inet_aton(ip))[0]
    prefix_int = int(prefix)
    mask = int('1' * prefix_int + '0' * (32 - prefix_int), 2)
    subnet_ip_int_min = ip_int & mask
    subnet_ip = socket.inet_ntoa(struct.pack("!I", subnet_ip_int_min))
    subnet_ip_int_max = subnet_ip_int_min | int('1' * (32 - prefix_int), 2)
    subnet_ip_max = socket.inet_ntoa(struct.pack("!I", subnet_ip_int_max))
    sys.stderr.write('Is %s between %s and %s\n' % (vip, subnet_ip, subnet_ip_max))
    if subnet_ip_int_min < vip_int < subnet_ip_int_max:
        subnet_ip = socket.inet_ntoa(struct.pack("!I", subnet_ip_int_min))
        print('%s/%s' % (subnet_ip, prefix))
        sys.exit(0)
sys.exit(1)
EOF
)
    ip -o addr show to "$net_cidr" | awk '{print $2}'
}

mkdir --parents /etc/keepalived

KEEPALIVED_IMAGE=registry.access.redhat.com/rhosp14/openstack-keepalived:14.0
if ! podman inspect "$KEEPALIVED_IMAGE" &>/dev/null; then
    echo "Pulling release image..."
    podman pull "$KEEPALIVED_IMAGE"
fi

API_DNS=$(sudo awk -F[/:] '/apiServerURL/ {print $5}' /opt/openshift/manifests/cluster-infrastructure-02-config.yml)
export MASTER_VIP=$(dig +noall +answer "$API_DNS" | awk '{print $NF}')
env INTERFACE=$(get_iface_in_vip_subnet "$MASTER_VIP") envsubst < /etc/keepalived/keepalived.conf.tmpl > /etc/keepalived/keepalived.conf

podman run \
        --rm \
        --volume /etc/keepalived:/etc/keepalived:z \
        --network=host \
        --cap-add=NET_ADMIN \
        "${KEEPALIVED_IMAGE}" \
        /usr/sbin/keepalived -f /etc/keepalived/keepalived.conf --dont-fork -D -l -P

# Workaround for https://github.com/opencontainers/runc/pull/1807
touch /etc/keepalived/.keepalived.done
