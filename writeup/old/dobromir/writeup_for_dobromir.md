# Risk Structured SIR Model

## Background

For Andrew Berdahl's agent-based modeling class I looked at epidemic dynamics
with heterogeneous risk with additional spatial/social structure.

I started with an agent-based SIR model: some number of "agents" interact
and spread a disease stochastically moving between susceptible,
infected and recovered states (each agent contacts every other
agent daily; infected agents spread disease to susceptible agents with
probability $\beta_c$; infected agents recover with probability $\gamma$)

Additionally, on a given day each individual in the population has
probability $p$ of choosing "risky" behavior (i.e. going to church or to a bar),
which is fixed per individual over the life of the simulation.

Agents who choose the "risky" behavior on a given day 
spread the disease to one-another with probability $\beta_r$ ($\beta_r > \beta_c$).



## Questions that I have

Before I started this project I read a few papers about super spreading events
(SSEs) and heterogeneity in epidemics:

https://journals.plos.org/plosbiology/article?id=10.1371/journal.pbio.3000897


https://journals.plos.org/plosone/article?id=10.1371/journal.pone.0250050


My main goal was to investigate a simple and tractable mechanism 
for heterogeneity and consider questions/themes that came up in these papers
concretely.

The two main questions:

### 1. How does heterogeneity affect the probability of the disease going extinct?

Basic theory of SSEs says that the more heterogeneous the population the less
likely a disease is to initially take off. I figured I could come up with any
of a few ways to quantify that in this model (analytically and in my simulation). 


### 2. Basic question of intervention timing

Suppose we have access to an intervention
measure like closing bars and restaurants, which in this model we could program
as setting $\beta_r$ to 0.

Since we may assume individuals with high riskyness (high $p$) are more likely to get
infected earlier on, is the timing of this intervention extremely sensitive?
Can I construct a scenario where an early intervention can prevent an epidemic
but a slightly later intervention has little or no impact; and make some statements
about the robustness of this scenario?

This is similar to questions raised in this paper:
but could help investigate the actual mechanics involved in one kind of SSE.

This is similar to this paper in preprint I worked on:

https://www.medrxiv.org/content/10.1101/2020.08.21.20179473v2

## DIFEQ model

I didn't actually think of how to model this deterministically for a while,
but I eventually came up with the exact same integro-differential model that you
and Mark Kot analyzed in your 2020 paper (in the continuous limit):


$$
\begin{aligned}
\frac{\partial S(p, t)}{\partial t} &=
	-\beta_c S(p, t) \int_{0}^1 I(u, t) du
	-\beta_r S(p, t) p \int_{0}^1 I(u, t) u du\\
\frac{\partial I(p, t)}{\partial t} &=
	\beta_c S(p, t) \int_{0}^1 I(u, t) du
	+ \beta_r S(p, t) p \int_{0}^1 I(u, t) u du - \gamma I(p, t)
\end{aligned}
$$

Here S(p, t) and I(p, t) are the populations of succeptible and infected
individuals at time $t$ with riskyness parameter $p$. $\beta_c$ and $\beta_r$
are paramaters that govern community and "risky" spread. And $u$ is the
dummy variable over which I'm taking integrals.

I'm not actually sure that this leads anywhere, but I noticed that the
integrals that came up looked like moments (as in probability) over the
riskyness (p) dimension, i.e.:

$\int_{0}^1 I(u, t) u du$ is the first moment of $I(p, t)$, and $\int_{0}^1 I(u, t) du$
is the zeroth moment.

Let $\bar I$ be the zeroth moment of $I(t, p)$ over $p$; $\hat I$ the first moment,
$\hat{\hat I}$ the second moment, etc; and the same for $S$.

Then:

$$
\begin{aligned}
\frac{\partial S}{\partial t} &=
	-\beta_c S \bar I
	-\beta_r S \hat I\\
\frac{\partial I}{\partial t} &=
	\beta_c S \bar I
	+ \beta_r S \hat I - \gamma I
\end{aligned}
$$

Now consider
$$
\begin{aligned}
\frac{d\hat I}{dt}
&= \frac{d}{dt}\left (\int_{0}^1 I(u, t) u du \right)\\
&= \int_{0}^1 \frac{\partial}{\partial t}I(u, t) u du\\
&= \int_{0}^1 (\beta_c S \bar I + \beta_r S u \hat I - \gamma I) u du\\
&= \beta_c \hat S \bar I + \beta_r \hat{\hat S} \hat I - \gamma \hat I\\
\end{aligned}
$$

and similarly

$$
\begin{aligned}
\frac{d\hat S}{dt}
&= \frac{d}{dt}\left (\int_{0}^1 S(u, t) u du \right)\\
&= \int_{0}^1 \frac{\partial}{\partial t}S(u, t) u du\\
&= \int_{0}^1 (-\beta_c S \bar I - \beta_r S u \hat I) u du\\
&= -\beta_c \hat S \bar I - \beta_r \hat{\hat S} \hat I\\
\end{aligned}
$$

And in general the derivative with respect to $t$ of any moment of $I$ or $S$ can be
expressed in terms of other moments of $I$ and $S$. This allows us to consider
the moment generating function $G_S(v, t) of S$:

$$G_S(v, t) = \bar S + u \hat S + \frac{u^2}{2!} \hat{\hat S} + ...$$

and its derivative:

$$
\begin{aligned}
\frac{\partial G_S}{\partial t}
	&= \frac{\partial \bar S}{\partial t} + u \frac{\partial \hat S}{\partial t}
	+ \frac{u^2}{2!} \frac{\partial \hat{\hat S}}{\partial t} + ...\\
	&= (-\beta_c \bar S \bar I - \beta_r \hat S \hat I) 
	+ u (-\beta_c \hat S \bar I - \beta_r \hat{\hat S} \hat I)
	+ \frac{u^2}{2!} (-\beta_c \hat{\hat S} \bar I - \beta_r \hat{\hat{\hat S}} \hat I) + ...\\
	&= -\beta_c \bar I (\bar S + u \hat S + \frac{u^2}{2!} \hat{\hat S} + ...)
	- \beta_r \hat I (\hat S + u \hat{\hat S} + \frac{u^2}{2!} \hat{\hat{\hat S}} + ...)\\
	&= -\beta_c \bar I G_S - \beta_r \hat I \frac{\partial G_S}{\partial u}
\end{aligned}
$$


A similar result holds for $G_I$, and you can write equations for both
generating functions:


$$
\begin{aligned}
\frac{\partial G_S(u, t)}{\partial t} &= -\beta_c \bar{I} G_S - \beta_r \hat{I} \frac{\partial}{\partial u} G_S\\
\frac{\partial G_I(u, t)}{\partial t} &= \beta_c \bar{I} G_S + \beta_r \hat{I} \frac{\partial}{\partial u} G_S - \gamma G_I
\end{aligned}
$$

So a third question I have:

### 3. Is there any interesting result I can get to using these generating function equations or further analyzing this model analytically?
