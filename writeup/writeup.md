# Writeup

## Outline

- Extinction Probability
- Peak Infecteds Size
- Final Size
- Evolution over time
- Timing of Intervention


## Extinction Probability

Using a branching process as an approximation

$\tau = G_X(\tau)$

We're interested in $E[\tau_i]$ for individual $i$ with riskiness $\rho_i$.

$$G_X(s) = P(X=0) + P(X=1)s^1 + P(X=2)s^2 + \ldots$$

$$P(X=0) = \rho_i \prod_j (1 - \rho_j \alpha_r) \prod_j (1 - \alpha_c) + (1 - \rho_i) \prod_j (1 - \alpha_c)$$

$$
\begin{aligned}
E[P(X=0)] &=
E[\rho_i \prod_j (1 - \rho_j \alpha_r) \prod_j (1 - \alpha_c) + (1 - \rho_i) \prod_j (1 - \alpha_c)]\\
&= E[\rho_i] \prod_j (1 - E[\rho_j] \alpha_r) \prod_j (1 - \alpha_c)
	* (1 - E[\rho_i]) \prod_j (1 - \alpha_c)\\
\end{aligned}
$$

The same result will hold for all $P(X=x)$, so
$E[G_X(\tau_i)] = E[\tau_i]$ depends only on $E[\rho]$