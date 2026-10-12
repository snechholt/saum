Note: i'm reading your answer and making notes incrementally. So thoughts are developing along/before you've possibly materialized the same thoughts later in your answer



On the data being returned (new state, findings, evidence). Who does the eval of this. It must be saum. So saum kind of acts as an expert on the project. Like if the test harness project did experiments (btw: why would it? From an internal standpoint that seems wasteful if there's a clear way to just achieve the result), then would it itself eval all the quirks and details of the experiment? If so - it won't know what is relevant to the architecture level. Would it just return an objective summary of what happened? If so, it risk accidentally drop details. E.g "experiment 1 spent most of its time creating an OCI registry and failed" when the crucial detail was that it did it in a "stupid" way but that was because of unclear specs. It's almost like saum should have access to everything that happened including commits and logs from every run. Then it can determine itself (or its human operator) how deep it wants to go. Want to just accept the summaries from the project? Or want to burn a lot of tokens analyzing what happened, and put it in perspective vs based on the project it derived from. Saum doesn't know about the internal details, methodology etc of the lower level project but can see it from the raw data, and that should be enough. Possibly we could force the projects to provide hints so saum doesn't need to read every line e.g. if it uses a subagent to run cmd commands, then that's less important to read in full vs the agent that managed that cmd agent. This does add pressure on saum to actually understand basically all projects (in different combinations). But that should be OK? Projects handle methodology internally but the end result is just data. Saum can read that data and add logic on top of it. It can contain the logic for "X and Y was considered in project 1, project A. Project A derives from project B. It seems B intended for X to be an option only but project A found a loophole. Suggestion to human: remove loophole in Project B, give context on how it was exploited". Project B would then internally be tasked with finding out *\_how\_* to fill the loophole. Saum doesn't say to do step 1,2,3 and it'll be fixed. It basically just reports a bug to project B. This workflow does demand a lot of saum. But is there a better way?

\--> On the decision to analyze-and-burn-tokens or not. Human would probably have triggered a "you should experiment" unless that is a native way for the project to act (if so, saum should know this so it can evaluate). Human can then decide if they want to back out of the experiment (too costly to follow up, even if it was explicitly asked for) OR: analysis could be delayed. Wait until we have budget to run it. It's just based on the persistent state (git commits, artifacts, logs etc), which is still there. Then do a follow-up check to see if any conclusions we come up with have already been fixed (e.g. loophole was fixed independently).



"Saum performs global impact analysis":

yes seems right. Otherwise we'd have to tell all the projects internally that "here is the spec. You've implemented everything up to A12, now we're at A13 and you need to figure out what's changed, then implement whatever is required". That would a) add this static logic to each project b) break with whatever methodology was inside the project (e.g. SpecKit as the anti-case of handling large changes (as an example, hypothetically we could have a tool that would just reject tasks as "too large, can't do it")). This last double parenthesis is actually interesting because maybe (just maybe - feel free to critique) the project should define some sort of "optimal task size" that saum would target. It could be a static rule and/or an interface between the two where they literally talked to each other while saum figured out the right decomposition.

\--> Again this will require a lot of saum.

\--> As with the philosophy of the previous paragraph of "saum just adds a bug report", saum needs to keep at the right level here as well. Not go into methodology. Just one or more issue tickets that the project can choose to interpret internally.



"Each project performs local impact analysis"

\`\`\`

Reconcile Foo\@F17 against Architecture\@A13.

Relevant change: Baz integration.

\`\`\`

Yes. But also: there might be many things in "Relevant change". That's the "problem". We don't have separate projects for "Foo, Export, Timeline and API". We have one "Implementation" project. And all your examples go into that one project. Like... maybe they would be different projects if they were each individual microservices and that would be a natural split. But another org may just set up one "implementation". Or "design" and "implementation". Basically whatever combo, saum needs to dynamically figure out what fits in there + ideally have read guidelines / communicated with the project to figure out what is an appropriate description to be taken over by the local impact analysis.



"Scenario 01's test-harness project might receive a very broad first reconciliation/materialization request and decide: There are four foundational uncertainties here. This is not ordinary implementation." How would it decide this, though? Then it needs to know about experiments and when to trigger them. Following on from the idea of direct communication between saum and project, we could add another endpoint where saum can ask "Is this a big task for you?" or "Do you think there's a lot of different approaches here" or similar. If yes (radical idea incoming), saum could basically *\_branch the whole project\_* and run experiments on each isolated instance of the project. I'm not sure on the followup: but one of either saum or the project then need to take all the info from that experiment and analyze it, then pick the best parts (that was the whole point of the experiments) and either experiment again or work on an implementation based on it.

\--> There's a discrepancy here in that the project gets not just one but multiple work items (like I said above). We might want to experiment on just one of those. So maybe when saum asks project if something is eligible for experiment, it is on a per-task basis and not on the whole thing?

\--> Did we talk about milestones in this context? In short: original though was saum needs to break down work into parts for projects to start working. This was at the point of my initial message, so things have changed a bit. But is it per milestone? I think there's not a 1-to-1 link between the old milestones concept and our new approach of "saum breaks down what A12->A13 means" or... is this actually what we meant? Either way, could be a good idea, based on how the current state is. I guess one thing about the milestones were that they were sequential. The handoff we give to the projects is not necessarily sequential - they could be implemented in any order and could also be merged. Is this also part of the saum-to-project communication during saum impact analysis? Not just settle what needs to do but define them as sequential milestones? We must be careful not to grow this interface too large but also consider what benefits we could get.



"Crucially, part of the experiment capability's contract is: report upstream-relevant findings, even when they didn't prevent local success." I just think this is hard to enforce. Hard to nail down the prompt to ensure nothing is lost. I've tried writing procedure-descibing skills before and have mostly failed if they involve more than just a few details. If that can be improved, sure, but if we could switch it around like I proposed before, then we've avoided the whole risk.



"A school strategy project could instead use:" interesting aspect when focusing on non software dev. But I'd counter here: the "alternatives" step is basically experiments. Then there's a workshop based on those experiments that chooses the best result based on that. That fits: saum starts experiments on isolated project "worktrees", then summarize differences. The workshop uses that, then comes back with a decision to saum. I'm just thinking how this would work with the proposed flow (saum triggers experiments directly). Yes, it could also be an internal process that saum is completely unaware of *\_and that is still supported\_* but in this case I don't see why it would need to then be doubly implemented by the project, unless the project uses another approach when producing alternatives. If same process, then having saum doing it basically feels like doing 2xA+2xB = 2(A+B), ie reducing the problem into simpler step that produces same result.



"Saum shouldn't necessarily author the final global plan from A13 by itself." yes I see we reached the same idea on saum-to-project communication. Which I guess is a sign that's a good idea. We don't need to focus too much on how direct/real-time the interface is unless it's relevant to another topic. But i think it would be nice as a loop on the comms and not just a sequence where saum just gets one round of input from the projects, then need to settle on the best design. Maybe it would be enough with just one data point, and I'm not sure what defines the exit condition for the loop, but if it's a complex breakdown, multiple rounds could be useful. And the question isn't really "here's what this means for foo" but rather "here's what this means for my methodology". Saum knows mostly about what it means for foo (right?). What it doesn't know is what this would mean for the way to handle Foo.



Also, the examples are a bit simplistic because you have done the breakdown into project-per-issue instead of the larger more realistic and dirty state. This is what I'm planning for in one of my projects

\`\`\`

architecture

  test-harness (end to end tests based on scenarios in architecture)

  design (system breakdown of architecture into components, interactions etc)

    implementation (code)

\`\`\`

implementation depends on \`system design\` AND \`test-harness\`, as it needs to conform to design and also pass tests defined in the harness. design and harness depends on architecture. In our case \`Baz\` would first affect \`test-harness\` and \`design\` in parallel. Then once they are both settled, \`implementation\` could be reconciled. With the new frameworks I've discovered recenly, maybe \`design\` has become redundant and a framework could work better with both design and implementation, but let's pretend this is the setup for future examples. A12->A13 would eventually mean that \`implementation\` gets a big batch of stuff that needs to be updated as a consequence.

\--> I see your "Don't create Saum nodes for workflow steps. Create them for independently authoritative state." and acknowledge



Addition to the above: scenario 01 in our case was for \`Traust\`, an app deployer/orchestrator, where scenario 01 tested the happy path including reading persistent state, fetching remote docker image from OCI registry, deploy container, plus more. The scenarios were open-ended with little direction on concrete implementation on how to test this. Tests need to create a fake environment of container registry etc for the deployer to work within. We found out each of these steps would ideally have a separate, focused round because the agents interpreted them wildly diffrerently. But that difference uncovered issues in the scenarios and architecture. Hence the value of experiments on open-ended tasks.



Whacky idea related to saum controlling experiments: if we had that structure, we could have an experiment on *\_all downstream projects\_*. Go wild and see how stuff turned out within each project and as an end result. Maybe we'd learn something. If tokens aren't a problem or the diff was small, it might be interesting. Not something to prioritize but it would be possible.



Counter arguments to saum controller experiments: it expands the surface area of saum, and in a slightly non-related way. We're handling a set of artefacts. Now suddenly we incorporate a workflow. Would the elevator pitch on the project contain a note on experiments?



From the towards-the-end of you reply: One thing I'm not sure if I understand/agree with is:

\`\`\`

findings ↑

evidence ↑

change proposals ↑

\`\`\`

Instinctively this feels like a risk in that the project may misinterpret the prompt and give garbage. We should find out how to mitigate (or avoid) that risk. The "can't implement, spec is self contradictory" is OK. Larger findings seem more risky.



Think critically about my suggestions. Don't just go along. I may be putting too much into saum.