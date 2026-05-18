from mininet.node import RemoteController, OVSSwitch
from mininet.cli import CLI
from mininet.log import setLogLevel, info
from mininet.link import TCLink

import sys
sys.path.append('/opt/containernet')
from mininet.net import Containernet


def myNetwork():
    net = Containernet(controller=RemoteController)

    info('*** Adding controller\n')
    # Подключаемся к Faucet, который слушает на хосте
    net.addController('c0', controller=RemoteController, ip='127.0.0.1', port=6653)

    info('*** Adding switches\n')
    # Используем OVS в режиме OpenFlow 1.3
    s1 = net.addSwitch('s1', cls=OVSSwitch, protocols='OpenFlow13')

    info('*** Adding docker containers\n')
    d1 = net.addDocker('d1', ip='10.0.0.1', dimage="service1:latest", mac='00:00:00:00:00:01')
    d2 = net.addDocker('d2', ip='10.0.0.2', dimage="service2:latest", mac='00:00:00:00:00:02')
    d3 = net.addDocker('d3', ip='10.0.0.3', dimage="service3:latest", mac='00:00:00:00:00:03')

    # Добавим "клиента", который будет генерировать нагрузку
    client = net.addDocker('client', ip='10.0.0.10', dimage="client:latest", mac='00:00:00:00:00:10')

    info('*** Creating links\n')
    # Соединяем всё со свитчем
    net.addLink(d1, s1, port1=1)
    net.addLink(d2, s1, port1=2)
    net.addLink(d3, s1, port1=3)
    net.addLink(client, s1, port1=4)

    info('*** Starting network\n')
    net.start()

    intf_name1 = d1.defaultIntf().name
    info(intf_name1)


    info('*** Configuring hosts\n')
    d1.cmd(f'ip addr add 10.0.0.100/32 dev {d1.defaultIntf().name}')
    d2.cmd(f'ip addr add 10.0.0.100/32 dev {d2.defaultIntf().name}')
    d3.cmd(f'ip addr add 10.0.0.100/32 dev {d3.defaultIntf().name}')
    client.cmd('arp -s 10.0.0.100 00:00:00:00:00:02')

    info('*** Running main processes\n')
    d1.cmd('./main &')
    d2.cmd('./main &')
    d3.cmd('./main &')

    info('*** Running CLI\n')
    CLI(net)

    info('*** Stopping network\n')
    net.stop()

if __name__ == '__main__':
    setLogLevel('info')
    myNetwork()