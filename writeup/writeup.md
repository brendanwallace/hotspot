# Writeup

## Outline

0. Abstract
1. Introduction
2. Superspreading
3. Extinction Probability
4. Max I, Total R
5. Disease Evolution
6. Intervention
7. Discussion

## 1. Introduction

[2005 Nature] outlines the importance of superspreading in disease transmission.
but only looks at individual variation of secondary infectiousness, no risk-taking
or any structure.
brings up the 80/20 rule (80% of infections caused by 20% of the population)

"For outbreaks avoiding stochastic extinction, epidemic growth rates strongly depend on variation in ν (Fig. 2c and Supplementary Fig. 2e, f). Diseases with high individual variation show infrequent but explosive epidemics after introduction of a single case. This pattern recalls SARS in 2003, for which many settings experienced no epidemic despite unprotected exposure to SARS cases27,28, whereas a few cities suffered explosive outbreaks8,9,10,15,26. Our results, using k̂ = 0.16 for SARS, explain this simply by the presence or absence of high-ν individuals in the early generations of each outbreak6. In contrast, conventional models (with k = 1 or k → ∞) cannot simultaneously generate frequent failed invasions and rapid growth rates without additional, subjective model structure."


[2020 laurent] argues that SSE are an important part of COVID spread.
[2021 dyani lewis] superspreading summary in COVID.
[2021 grossman] argues that heterogeneity is important in understanding COVID.

Here we introduce a model for disease transmission which allows for the close
examination of location-based SSE. We use this model to investigate
how this mechanism affects:
- the likelihood of an epidemic outbreak upon introduction to a small subpopulation
- the peak (largest number of simultaneously infected inviduals) and final
(recovered individuals at the end) size of the outbreak
- the efficacy of different types of preventative and control interventions.

We find that presence of this heterogeneity makes an outbreak less likely in
nearly all cases. When there is an outbreak, increasing heterogeneity
accelerates both the early growth of the disease and its later decay. When
average infectiousness is low (R0 < 1.2) this causes a higher peak and
larger final size, when average infectiousness is moderate (R0 1.2-2.0), it causes
a higher peak but smaller final size of outbreaks, and when infectiousness is
high (R0 > 2.0) it causes a lower peak and final size.

[e.g.: Interventions targetted at the high-risk segment are time sensitive, 
fluctuating risk tolerance could help explain why it doesn't ever go away]


## 2. Extinction Probability

Using a branching process as an approximation

$\tau = G_X(\tau)$

We're interested in $E[\tau_i]$ for individual $i$ with riskiness $\rho_i$.

$$G_X(s) = P(X=0) + P(X=1)s^1 + P(X=2)s^2 + \ldots$$

$$P(X=0) = \rho_i \prod_j (1 - \rho_j \alpha_r) \prod_j (1 - \alpha_c) + (1 - \rho_i) \prod_j (1 - \alpha_c)$$

$$
\begin{aligned}
E[P(X=0)] &=
E[\rho_i] \prod_j (1 - \rho_j \alpha_r) \prod_j (1 - \alpha_c) + (1 - \rho_i) \prod_j (1 - \alpha_c)]\\
&= E[\rho_i] \prod_j (1 - E[\rho_j] \alpha_r) \prod_j (1 - \alpha_c)
	* (1 - E[\rho_i]) \prod_j (1 - \alpha_c)\\
\end{aligned}
$$

The same result will hold for all $P(X=x)$, so
$E[G_X(\tau_i)] = E[\tau_i]$ depends only on $E[\rho]$