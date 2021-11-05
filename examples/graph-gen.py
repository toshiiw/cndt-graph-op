import getopt
import random
import sys

class Num2Vertex:
    def __init__(self):
        self.m = {}

    def v(self, i):
        try:
            return self.m[i]
        except KeyError:
            pass

        if i < 26:
            self.m[i] = chr(ord('a') + i)
            return self.m[i]
        else:
            il = []
            io = i
            while i >= 26:
                il.append(self.v(i % 26))
                i //= 26
            il.append(self.v(i))
            il.reverse()
            self.m[io] = ''.join(il)
            return self.m[io]

def f(name, vertices, weight, directed=False):
    print("""apiVersion: graph.example.valinux.co.jp/v1alpha1
kind: %sEdge
metadata:
  name: %s
spec:""" % ("Di" if directed else "", name))
    print(("  from: \"%s\"\n  to: \"%s\"\n  capacity: %s" if directed else "  vertex:\n  - \"%s\"\n  - \"%s\"\n  weight: %s") %
            (vertices[0], vertices[1], weight))
    print("---")

def g(sz, directed=False):
    n2v = Num2Vertex()
    vi = 0
    for i in range(sz):
        cur = n2v.v(i)
        for distance, prob, flip in [(3, 1, 0), (2, .5, .1), (6, .3, .5)]:
            if i - distance >= 0:
                if prob == 1 or random.random() < prob:
                    v2 = [n2v.v(i - distance), cur]
                    if random.random() < flip:
                        v2.reverse()
                    f("e%d" % vi, v2, random.randint(10,20), directed)
                    vi += 1
    return n2v.v(sz - 1)

def generate_problem(lv, directed):
    print("""apiVersion: graph.example.valinux.co.jp/v1alpha1
kind: %s
metadata:
  name: problem
spec:""" % ('MaxFlow' if directed else 'MaxCut'))
             
    if directed:
        print("  from: a\n  to: \"%s\"" % lv)
    else:
        print("  iteration: 100")

if __name__ == '__main__':
    opts, args = getopt.getopt(sys.argv[1:], 'dhp')
    directed = False
    gen_prob = False
    for o, a in opts:
        if o == '-d':
            directed = True
        elif o == '-p':
            gen_prob = True
        else:
            print("""Usage: %s [-d] [-p] <num of vertices>
\t-d: Generate a directed graph
\t-p: Also generate a maxflow/maxcut problem""" % sys.argv[0])
            sys.exit(0)
    last_v = g(int(args[0]), directed)
    if gen_prob:
        generate_problem(last_v, directed)
