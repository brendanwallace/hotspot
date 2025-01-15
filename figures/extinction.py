import numpy as np

import settings

N = 1000

START = 0.0
END = settings.R0_END
STEP = 0.1


def binomial_G(g, beta_c, beta_h, mean_p, N):
    a, b, p = beta_c, beta_h, mean_p # makes the equation a little easier to write
    return p*(((1 - p*b)*(1 - a) + (1 - (1-p*b)*(1-a))*g)**N) + (1 - p)*((1 - a + a * g)**N)

def poisson_G(g, beta_c, beta_h, mean_p, N):
    a, b, p = beta_c, beta_h, mean_p # makes the equation a little easier to write
    return (1 - p)*(np.exp(N*a*(g - 1))) + p*(np.exp(N*(a + b * p)*(g - 1)))

def theoretical_extinction_binomial(beta_c, beta_h, mean_p, G, N=1000):
    """
    Computes theoretical extinction probability.
    """
    upper_g = 1
    lower_g = 0
    for T in range(1000):
        g = (upper_g + lower_g)/2
        g_ = G(g, beta_c, beta_h, mean_p, N)
        if g_ < g:
            upper_g = g
        else:
            lower_g = g
    return g


def theoretical_extinction_poisson(beta_c, beta_h, mean_p, N=1000):
    a, b, p = beta_c, beta_h, mean_p # makes the equation a little easier to write
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


# homogeneous_outbreak, poisson_outbreak = [], []
# for R0 in X:
#     homogeneous_outbreak.append(1.0 - theoretical_extinction((R0)/N, 0, 0, binomial_G))
#     poisson_outbreak.append(1.0 - poisson_extinction(R0))


X = np.arange(START, END + STEP, STEP)
HOMOGENEOUS = []
for R0 in X:
	HOMOGENEOUS.append(1.0 - poisson_extinction(R0))