from random import random, seed, randrange
seed(11)
ativos = [] # onde cada elemento tem duas componentes: indicador de lucro esperado e risco

def ordenar_ativos(e):
  return e[0]/e[1]

for i in range(20):
    ativos.append([random(), random()])
ativos.sort(key=ordenar_ativos)

##################################################################################################

cardinalidade = 10
montante = 100000.0

restricao = 0.2
min_restricao = 0.1
vizinhanca = 100
offset = 100

L = range(cardinalidade) 
s = [montante/cardinalidade] * cardinalidade

def incrementar(lista, indice, valor):
    restante_anterior = montante - lista[indice]
    if lista[indice] + valor >= montante * restricao:
        lista[indice] = montante * restricao
    else:
        lista[indice] = lista[indice] + valor
    
    restante = montante - lista[indice]

    for i in range(cardinalidade):
        if i != indice:
            lista[i] = (float(lista[i])/float(restante_anterior)) * float(restante)

    
    return lista

def fitness(s):
    fit = 0

    for i in range(cardinalidade):
        fit += (float(s[i])/float(montante)) * L[i]

    return fit

def melhor_vizinho(s, k):
    best = list(s)
    best_fit = fitness(best)
    i = 0
    while(i < cardinalidade):
        viz_ = list(best)
        for j in range(1, k+1):

            if float(viz_[i]) + offset*j <= restricao*montante and float(viz_[i]) + offset*j >= min_restricao*montante:
                viz_ = incrementar(viz_, i, offset*j)
                fit = fitness(viz_)
                if fit > best_fit:
                    best_fit = fit
                    best = list(viz_)
                    i = 0
        i+= 1

    return best, best_fit

def vnd(init):
    best_fit = fitness(init)
    best = list(init)
    for k in range(vizinhanca):
        b, b_fit = melhor_vizinho(best, k+1)
        if b_fit > best_fit:
            best_fit = b_fit
            best = b
    return best, best_fit


print('Fitness S ', fitness(s))
best, best_fit = vnd(s)
#print('VND ', best_fit)
print('VND ', best)

########################################## VNS


def mutate(_s):
    s_copy = list(_s)
    s_i = randrange(0, cardinalidade)

    incrementar(s_copy, s_i, montante * 0.05) # pega um vizinho aleatorio

    return s_copy


vns_run_times = 20
s_vns = list(s)
best = list(s)
best_fit = 0

for i in range(vns_run_times):
    b, b_fit = vnd(s_vns)
    if b_fit > best_fit:
        best_fit = b_fit
        best = b

    s_vns = mutate(best)

#print('VNS ', best_fit)
print('VNS ', best)
