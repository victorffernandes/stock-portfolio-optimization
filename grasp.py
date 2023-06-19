from random import random, seed, randrange

ativos = []


cardinalidade = 20
max_grasp = 20
seed(10)

s_ = -1

def ordenar_ativos(e):
  return e[0]/e[1]

for i in range(100):
    ativos.append([random(), random()])
ativos.sort(key=ordenar_ativos)

# def get_elementos_lrc(lrc, lista, valor, visited):
#     if valor == -1:
#         return lista[-lrc:]
#     else:
#         lista_lrc = []
#         for x in range(lista):



def grasp_construcao(lrc):
    carteira = [-1] * cardinalidade
    media = 0

    for i in range(cardinalidade):

        #lr = list([(ativos[i][0]/ativos[i][1]), (ativos[i][0]/ativos[i][1]) - media, i] for i in range(len(ativos)) if i not in carteira)
        lr = list([(ativos[i][1]), media - ativos[i][1], i] for i in range(len(ativos)) if i not in carteira)
        lr.sort(key=lambda x: x[1])
        rand_chosen = randrange(0, 3)


        # while()

        carteira[i] = lr[-lrc:][rand_chosen][2]

        #media += lr[-lrc:][rand_chosen][0] / cardinalidade
        media += lr[-lrc:][rand_chosen][0]
        media /= i + 1
    carteira.sort()
    return carteira
    #print('hey ', lr)

        # ativos[-lrc]
        
    # for i in 


# for x in range(max_grasp):
#     s = grasp_construcao(3)

solutions = {}

for g in range(max_grasp):
    r = hash(frozenset(grasp_construcao(3)))
    rs = str(r)
    if solutions.get(rs):
        print('equal', rs)
    else:
        print('not equal', rs)
        solutions[rs] = 1


# for i in range(len(ativos)):
#     print(ativos[i][0]/ativos[i][1], i)

