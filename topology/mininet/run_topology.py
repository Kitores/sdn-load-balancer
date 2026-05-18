#!/usr/bin/python3
from mininet.net import Containernet
from mininet.node import RemoteController, OVSSwitch
from mininet.cli import CLI
from mininet.log import setLogLevel, info
from mininet.link import TCLink

def myNetwork():
    net = Containernet(controller=RemoteController)

    info('*** Adding controller\n')
    # Подключаемся к Faucet, который слушает на хосте
    net.addController('c0', controller=RemoteController, ip='127.0.0.1', port=6653)

    info('*** Adding switches\n')
    # Используем OVS в режиме OpenFlow 1.3
    s1 = net.addSwitch('s1', cls=OVSSwitch, protocols='OpenFlow13')

    info('*** Adding docker containers\n')
    # Твои Go-сервисы
    d1 = net.addDocker('d1', ip='10.0.0.1', dimage="service1:latest")
    d2 = net.addDocker('d2', ip='10.0.0.2', dimage="service2:latest")
    d3 = net.addDocker('d3', ip='10.0.0.3', dimage="service3:latest")
    
    # Добавим "клиента", который будет генерировать нагрузку
    # Можно использовать обычный образ ubuntu или alpine
    client = net.addDocker('client', ip='10.0.0.10', dimage="client:latest")

    info('*** Creating links\n')
    # Соединяем всё со свитчем
    net.addLink(d1, s1)
    net.addLink(d2, s1)
    net.addLink(d3, s1)
    net.addLink(client, s1)

    info('*** Starting network\n')
    net.start()

    
    info('*** Running CLI\n')
    CLI(net)

    info('*** Stopping network\n')
    net.stop()

if __name__ == '__main__':
    setLogLevel('info')
    myNetwork()