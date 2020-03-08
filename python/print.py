import networkx as nx
import numpy as np
import matplotlib.pyplot as plt
import pylab
import json

def color_choice(x):
  if x in qo and x in qf:
    return 'green'
  elif x in qo:
    return 'red'
  elif x in qf:
    return 'blue'
  else:
    return 'grey'


G = nx.DiGraph()

file1 = open('./python/graph.txt', 'r') 
Lines = file1.readlines()

edge_labels=dict()
edges = []

for line in Lines: 
  line = line.replace("\n", "")
  A, trans, B = line.split(",")
  edges.append((A, B))
  edge_labels[(A, B)] = trans

file2 = open('./python/important.txt', 'r') 
Lines = file2.readlines()
qo, qf = Lines
qo = qo.replace("[","").replace("]", "").replace("\n", "").split(" ")
qf = qf.replace("[","").replace("]", "").replace("\n", "").split(" ")

print(qo, qf)

print(edges)

G.add_edges_from(edges)

values = [node for node in G.nodes()]
values = list(map(color_choice, values))

edge_colors = ['black' for edge in G.edges()]

node_labels = {node:node for node in G.nodes()}
pos=nx.spring_layout(G)
nx.draw_networkx_labels(G, pos, labels=node_labels)
nx.draw_networkx_edge_labels(G,pos,edge_labels=edge_labels)
nx.draw(G,pos, node_color = values, node_size=100,edge_color=edge_colors,edge_cmap=plt.cm.Reds)
pylab.show()