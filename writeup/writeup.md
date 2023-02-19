---
Title: Problem-place model of risk structured disease spread
Author: Brendan Wallace
---

## Outline

0. Abstract
1. Introduction
2. Model
	a. agent based model, introduced textually
	b. branching process
	c. integro differential
	d. difference
3. Results
	a. Extinction Probability
	b. Max I, Total R
	c. Disease Evolution
4. Discussion

## 0. Abstract

This is the abstract.

## 1. Introduction

[2005 Nature] outlines the importance of superspreading in disease transmission.
but only looks at individual variation of secondary infectiousness, no risk-taking
or any structure.
brings up the 80/20 rule (80% of infections caused by 20% of the population)

"For outbreaks avoiding stochastic extinction, epidemic growth rates strongly depend on variation in ν (Fig. 2c and Supplementary Fig. 2e, f). Diseases with high individual variation show infrequent but explosive epidemics after introduction of a single case. This pattern recalls SARS in 2003, for which many settings experienced no epidemic despite unprotected exposure to SARS cases27,28, whereas a few cities suffered explosive outbreaks8,9,10,15,26. Our results, using k̂ = 0.16 for SARS, explain this simply by the presence or absence of high-ν individuals in the early generations of each outbreak6. In contrast, conventional models (with k = 1 or k → ∞) cannot simultaneously generate frequent failed invasions and rapid growth rates without additional, subjective model structure."


[2020 laurent] argues that SSE are an important part of COVID spread.
[2021 dyani lewis] superspreading summary in COVID.
[2021 grossman] argues that heterogeneity is important in understanding COVID.

[200
https://www.nature.com/articles/s41586-020-2923-3

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


\[ background - super spreading, why and how it's relevant here \]
Disease spread is not homogeneous – it's well understood now that during any
outbreak, some individuals (commonly referred to as super-spreaders) are
responsible for significantly more infections than others [citations].
During the COVID-19 pandemic in particular,  

\[ limitations of existing knowledge, how individual variance is insufficient
and empirically, problem places are important. \]


\[ problem place model to the rescue \]

Here, we present a simple model of a disease outbreak that nevertheless
captures risk taking behavior and the difference in how readily disease spreads
in certain "problem places", (bars, restaurants, churches etc.). We analyze
this model and find a few striking differences compared to a homogeneous SIR
model, _even when picking parameters to ensure an indentical basic reproduction
number_. We see:

- a smaller chance of an outbreak, when one agent is infected in a finite population
	(is there a good word for this?)
- faster, more "explosive" outbreak when they do occur,
- with mixed overall outcomes in terms of peak and total
	numbers of infections.

In the discussion section, we show how this mechanism is largely in agreement
with Laurants' findings about super spreader events. We look at some recent work
into patterns in urban mobility and discuss which "problem place" archetypes
are most relevent.



## 2. Model

### 2.1 Agent-based model

We primarily use an agent based model to investigate this scenario.

#### Setup

To run a simulation, we initialize a fixed-size population of $N$ agents.
Each agent $i$ ($i \in \[0, N\)$) is characterized by a fixed "riskyness" parameter
$\rho_i \in \[0, 1\]$ that never changes, and a disease state of Susceptible (S),
Infected (I) or Recovered (R). We set the following parameters:

- $\beta_c \in [0, 1], << 1$ community spread rate – the probability of disease spreading
from an I to an S individual through a single community contact,
-  $\beta_p \in [0, 1], << 1$ problem-place spread rate - the probability of disease spreading
from an I to an S individual through contact in the problem-place.
- $\gamma \in (0, 1]$ the recovery rate - the rate at which I individuals move to
R (the disease lasts on average $1/\gamma$ time units).

$\beta_p$ is typically much larger than $\beta_c$, though both are much
smaller than 1.

All agents $i$ are initially S, with $\rho_i$ drawn iid. from a fixed
distribution. One agent is chosen at random and set to I.


#### Dynamics

##### Risk Taking
Every timestep, each agent $i$ is set as visiting the problem place with
probability $\rho_i$.

##### Disease Spread
All agents in the problem place make a contact with all other agents in the
problem place; I individuals in this subset spread the disease to S
individuals with probability $\beta_p$ (problem place spread).
Simultaneously, all agents (regardless of whether they have visited the problem
place this time step) make contact with all other agents and I individuals
spread to S individuals with probability $\beta_c$ (homogeneous community
spread).

##### Recovery
I agents recover with probability $\gamma$. If there are no I agents, the
simulation ends.



### 2.2 Branching process model

To aid in investigating the dynamics of this model when the number of infected
agents is small; we use ....


[\this is developed more in the text\]



### 2.3 Integro-differential model

### 2.4 Difference model

(Introduced super briefly.)

### 2.5 R0

### Results

#### Outbreak Probability

The model is stochastic, so we see different sizes of outbreaks.

\[ diagram of homogeneous case on the left, 50/50 on the right. R0 = 2 \]

The majority of outbreaks are large, but some end up being small.


In the homogeneous case, we can nicely predict the probility of an outbreak
by modeling the process as a branching process. (explains branching process
model)

$$\psi = G(\psi)$$

This is plotted in (figure)

\[ diagram of predicted vs actual, R0 on the x axis, extinction on the y axis \]


We can apply the same process here:

(developes branching process model)

(presents extinction probability)

(graph of extinction probabilities vs predictied)

(brief discussion – the so what?)


#### Max I, Total R