#!/usr/bin/python3

import yaml
from scapy.all import *
import click

def yaml_read(yaml_file):
    with open(yaml_file, 'r') as f:
        data = yaml.load(f, Loader=yaml.FullLoader)
        packet = data['packet']
    return packet

def create_ether(eth):
    p = Ether()
    if 'dst'  in eth: p.dst  = eth['dst']
    if 'src'  in eth: p.src  = eth['src']
    if 'type' in eth: p.type = eth['type']
    return p

def create_ipv6(ipv6):
    p = IPv6()
    if 'version' in ipv6: p.version = ipv6['version']
    if 'tc'      in ipv6: p.tc      = ipv6['tc']
    if 'fl'      in ipv6: p.fl      = ipv6['fl']
    if 'plen'    in ipv6: p.plen    = ipv6['plen']
    if 'nh'      in ipv6: p.nh      = ipv6['nh']
    if 'hlim'    in ipv6: p.hlim    = ipv6['hlim']
    if 'src'     in ipv6: p.src     = ipv6['src']
    if 'dst'     in ipv6: p.dst     = ipv6['dst']
    return p

def create_srh(srh):
    p = IPv6ExtHdrSegmentRouting()
    if 'nh'          in srh: p.nh         = srh['nh']
    if 'len'         in srh: p.len        = srh['len']
    if 'type'        in srh: p.type       = srh['type']
    if 'segleft'     in srh: p.segleft    = srh['segleft']
    if 'lastentry'   in srh: p.lastentry  = srh['lastentry']
    if 'unused1'     in srh: p.unused1    = srh['unused1']
    if 'protected'   in srh: p.protected  = srh['protected']
    if 'oam'         in srh: p.oam        = srh['oam']
    if 'alert'       in srh: p.alert      = srh['alert']
    if 'hmac'        in srh: p.hmac       = srh['hmac']
    if 'unused2'     in srh: p.unused2    = srh['unused2']
    if 'tag'         in srh: p.tag        = srh['tag']
    if 'addresses'   in srh: p.addresses  = srh['addresses']
    if 'tlv_objects' in srh: p.tlv_objects= srh['tlv_objects']
    return p

def create_ip(ip):
    p = IP()
    if 'version' in ip: p.version = ip['version']
    if 'ihl'     in ip: p.ihl     = ip['ihl']
    if 'tos'     in ip: p.tos     = ip['tos']
    if 'len'     in ip: p.len     = ip['len']
    if 'id'      in ip: p.id      = ip['id']
    if 'flags'   in ip: p.flags   = ip['flags']
    if 'frag'    in ip: p.frag    = ip['frag']
    if 'ttl'     in ip: p.ttl     = ip['ttl']
    if 'proto'   in ip: p.proto   = ip['proto']
    if 'chksum'  in ip: p.chksum  = ip['chksum']
    if 'src'     in ip: p.src     = ip['src']
    if 'dst'     in ip: p.dst     = ip['dst']
    if 'options' in ip: p.options = ip['options']
    return p

def create_icmp(icmp):
    p = ICMP()
    if 'type'       in icmp: p.type       = icmp['type']
    if 'code'       in icmp: p.code       = icmp['code']
    if 'chksum'     in icmp: p.chksum     = icmp['chksum']
    if 'id'         in icmp: p.id         = icmp['id']
    if 'seq'        in icmp: p.seq        = icmp['seq']
    if 'ts_ori'     in icmp: p.ts_ori     = icmp['ts_ori']
    if 'ts_rx'      in icmp: p.ts_rx      = icmp['ts_rx']
    if 'ts_tx'      in icmp: p.ts_tx      = icmp['ts_tx']
    if 'gw'         in icmp: p.gw         = icmp['gw']
    if 'ptr'        in icmp: p.ptr        = icmp['ptr']
    if 'reserved'   in icmp: p.reserved   = icmp['reserved']
    if 'length'     in icmp: p.length     = icmp['length']
    if 'addr_mask'  in icmp: p.addr_mask  = icmp['addr_mask']
    if 'nexthopmtu' in icmp: p.nexthopmtu = icmp['nexthopmtu']
    return p

def create_udp(udp):
    p = UDP()
    if 'sport'  in udp: p.sport  = udp['sport']
    if 'dport'  in udp: p.dport  = udp['dport']
    if 'len'    in udp: p.len    = udp['len']
    if 'chksum' in udp: p.chksum = udp['chksum']
    return p

def create_tcp(tcp):
    p = TCP()
    if 'sport'    in tcp: p.sport    = tcp['sport']
    if 'dport'    in tcp: p.dport    = tcp['dport']
    if 'seq'      in tcp: p.seq      = tcp['seq']
    if 'ack'      in tcp: p.ack      = tcp['ack']
    if 'dataofs'  in tcp: p.dataofs  = tcp['dataofs']
    if 'reserved' in tcp: p.reserved = tcp['reserved']
    if 'flags'    in tcp: p.flags    = tcp['flags']
    if 'window'   in tcp: p.window   = tcp['window']
    if 'chksum'   in tcp: p.chksum   = tcp['chksum']
    if 'urgptr'   in tcp: p.urgptr   = tcp['urgptr']
    if 'options'  in tcp: p.options  = tcp['options']
    return p

def create_packet(p):
    packet = ""
    ### outer packet
    if 'ether' in p:
        packet /= create_ether(p['ether'])
    if 'ipv6' in p:
        packet /= create_ipv6(p['ipv6'])
    if 'srh' in p: # IPv6ExtHdrSegmentRouting
        packet /= create_srh(p['srh'])
    if 'ip' in p:
        packet /= create_ip(p['ip'])
    if 'icmp' in p:
        packet /= create_icmp(p['icmp'])
    if 'udp' in p:
        packet /= create_udp(p['udp'])
    if 'tcp' in p:
        packet /= create_tcp(p['tcp'])
    ### inner packet
    # We could call recursively by defining 'packet' and 'inner_packet' on same level.
    # However, keeping as below just in case we want to change order or available
    # header between outer and inner packet.
    if 'inner_packet' in p:
        i = p['inner_packet']
        if 'ether' in i:
            packet /= create_ether(i['ether'])
        if 'ipv6' in i:
            packet /= create_ipv6(i['ipv6'])
        if 'ip' in i:
            packet /= create_ip(i['ip'])
        if 'icmp' in i:
            packet /= create_icmp(i['icmp'])
        if 'udp' in i:
            packet /= create_udp(i['udp'])
        if 'tcp' in i:
            packet /= create_tcp(i['tcp'])

    return packet

class sendpacket(object):
    def __init__(self, yaml_file='packet.yaml'):
        self.send_iface = 'lo'
        self.yaml_file = yaml_file
        self.packet = yaml_read(self.yaml_file)
        self.send_iface = 'lo'
        if 'send_iface' in self.packet:
            self.send_iface = self.packet['send_iface']
    
    def send(self, p, c=0):
        sendp(p, iface=self.send_iface, count=c)

@click.command()
@click.argument('file_yaml')
@click.option('-c', '--count', default=1, help='number of packets to send')
@click.option('-i', '--send_iface', default='', help='interface name to send packet')
@click.option('-s', '--show', is_flag=True, help='show packet details')
@click.option('--debug', is_flag=True, help='show debug messages')
def cmd(file_yaml, count, send_iface, show, debug):
    sp = sendpacket(file_yaml)
    p = create_packet(sp.packet) #scapy packet

    if send_iface: sp.send_iface = send_iface
    if debug:
        print("DEBUG: sp.send_iface: {}".format(sp.send_iface))
        print("DEBUG: ", sp.packet)
    if show: p.show()

    try:
        socket.if_nametoindex(sp.send_iface)
    except OSError as e:
        print(e)
    else:
        sp.send(p, count)

if __name__ == '__main__':
    cmd()

