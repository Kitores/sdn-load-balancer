from diagrams import Diagram, Cluster, Edge
from diagrams.onprem.monitoring import Prometheus
from diagrams.onprem.client import User
from diagrams.onprem.container import Container
from diagrams.onprem.network import Faucet
from diagrams.providers.generic.network import Switch

# Стилизация для кластеров
cluster_attr = {
    "bgcolor": "#f4f7eb",
    "color": "#cbd5e1",
    "fontsize": "12"
}

graph_attr = {
    "splines": "spline",
    "bgcolor": "#e0f2fe",
    "nodesep": "0.6",
    "ranksep": "0.8"
}

with Diagram(name="Course Work: Сети ЭВМ", show=False, direction="TB", graph_attr=graph_attr, outformat="png"):
    
    with Cluster("Topology", graph_attr={"bgcolor": "#e3f2fd", "color": "#b0bec5"}):
        
        # 1. Кластер Prometheus
        with Cluster("Prometheus", graph_attr=cluster_attr):
            prometheus = Prometheus("\n\nhttp://localhost:9090")
            
        # 2. Кластер Load Emulator
        with Cluster("Load Emulator", graph_attr=cluster_attr):
            emulator = User("Emulator\nIP: 10.0.0.10")
            
        # 3. Кластер Service 2
        with Cluster("Service 2", graph_attr=cluster_attr):
            service2 = Container("Alpine Container\nIP: 10.0.0.2")
            
        # 4. Кластер Service 1
        with Cluster("Service 1", graph_attr=cluster_attr):
            service1 = Container("Alpine Container\nIP: 10.0.0.1")
            
        # 5. Кластер SDN-Controller
        with Cluster("SDN-Controller", graph_attr={"bgcolor": "#f1f5f9", "color": "#cbd5e1"}):
            faucet = Faucet("faucet OpenFlow1.3")
            ovs = Switch("Open vSwitch")
            
            faucet >> Edge(color="#64748b") >> ovs
        
        prometheus >> Edge(label="http://localhost:9302/metrics", color="#64748b") >> ovs
        
        # Двунаправленные сетевые линки от OVS к хостам
        ovs << Edge(color="#64748b") >> emulator
        ovs << Edge(color="#64748b") >> service2
        ovs << Edge(color="#64748b") >> service1