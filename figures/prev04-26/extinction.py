#!/usr/bin/env python
# coding: utf-8

# In[2]:


import numpy as np
import matplotlib.pyplot as plt

N = 1000


# In[25]:


def binomial_G(g, alpha_c, alpha_r, mean_p, N):
    a, b, p = alpha_c, alpha_r, mean_p # makes the equation a little easier to write
    return p*(((1 - p*b)*(1 - a) + (1 - (1-p*b)*(1-a))*g)**N) + (1 - p)*((1 - a + a * g)**N)

def poisson_G(g, alpha_c, alpha_r, mean_p, N):
    a, b, p = alpha_c, alpha_r, mean_p # makes the equation a little easier to write
    return (1 - p)*(np.exp(N*a*(g - 1))) + p*(np.exp(N*(a + b * p)*(g - 1)))

def theoretical_extinction(alpha_c, alpha_r, mean_p, G, N=1000):
    """
    Computes theoretical extinction probability.
    """
    upper_g = 1
    lower_g = 0
    for T in range(1000):
        g = (upper_g + lower_g)/2
        g_ = G(g, alpha_c, alpha_r, mean_p, N)
        if g_ < g:
            upper_g = g
        else:
            lower_g = g
    return g

start = 0.0
end = 8.0
step = 0.1
X = np.arange(start, end + step, step)

def theoretical_extinction_p(alpha_c, alpha_r, mean_p, N=1000):
    a, b, p = alpha_c, alpha_r, mean_p # makes the equation a little easier to write
    upper_g = 1
    lower_g = 0
    for T in range(1000):
        g = (upper_g + lower_g)/2
        g_ = p*(((1 - p*b)*(1 - a) + (1 - (1-p*b)*(1-a))*g)**N) + (1 - p)*((1 - a + a * g)**N)
        if g_ < g:
            upper_g = g
        else:
            lower_g = g
    return g

def poisson_extinction(R0):
    """
    Computes theoretical extinction probability.
    """
    upper_g = 1
    lower_g = 0
    for T in range(1000):
        g = (upper_g + lower_g)/2
        g_ = np.e ** (R0 * (g - 1))
        if g_ < g:
            upper_g = g
        else:
            lower_g = g
    return g


homogeneous_outbreak, poisson_outbreak = [], []
for R0 in X:
    homogeneous_outbreak.append(1.0 - theoretical_extinction((R0)/N, 0, 0, binomial_G))
    poisson_outbreak.append(1.0 - poisson_extinction(R0))


# In[10]:


for N in [5, 10, 100, 1000]:
    homogeneous_outbreak = []
    for R0 in X:
        homogeneous_outbreak.append(1.0 - theoretical_extinction((R0/N), 0, 0, binomial_G, N=N))
    plt.plot(X, homogeneous_outbreak, label=f"N={N}")
plt.plot(X, poisson_outbreak, label="poisson")
plt.legend()


# In[33]:


for N in [5, 10, 100, 1000]:
    for rf in [0.75]:
        for risk_mean in [0.125, 0.25, 0.5]: 
                B = []
                P = []
                for R0 in X:
                    alpha_c = R0 / N * (1-rf)
                    alpha_r = R0 / N * rf / risk_mean / risk_mean
                    B.append(1.0 - theoretical_extinction(alpha_c, alpha_r, risk_mean, binomial_G, N=N))
                    P.append(1.0 - theoretical_extinction(alpha_c, alpha_r, risk_mean, poisson_G, N=N))
                plt.plot(X, B, '--')
                plt.plot(X, P)
    plt.show()


# In[ ]:




