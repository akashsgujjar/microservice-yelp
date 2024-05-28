# Lab 4: Open-Ended Assignment (22.5 pts)

Lab 4 is an open-ended end-of-term project, to be done in the same lab
groups as previous assignments. In scoping the project, we encourage
lab groups to target a project that will take you roughly the same
amount of time as each of the other labs.

We provide a set of suggestions for possible topics for the lab, but
you are free to design a project of your own choosing. All projects
should be on some topic related to datacenter or cloud technologies,
have some implementation component (hardware or software), and produce
some measurable or demonstrable output.

## Deliverables

There are several deliverables (per group):

  * **Proposal**: approximately a half-page describing your plans for
    the project, including a description of the plans for measurable
    or demonstrable output. This is due prior to section labeled "Lab
    4 plan due (before section)" in the class calendar. Section will
    consist of groups exchanging plans for student peer review of
    everyone's plans.

  * **Project writeup**: two or three pages describing the project
    along with the measurable or demonstrable output, due by the
    deadline labeled "Lab 4 Report due" in the class calendar (note
    the exact time). The report should clearly identify which student
    worked on which part of the project.

  * **Tar file with code**. The project can be completed in any
    programming language you find convenient, with the tar file due at
    the same time as the project writeup.

## Project Suggestions

These are meant as suggestions; student groups are free to design your
own project on any topic related to the class.

Collaboration between groups, e.g. to tackle a deeper problem in
aggregate than any single group can do on their own, is allowed. For
this, we expect a longer combined report (e.g., for two groups, the
writeup can be twice as long). Please try to divide the work along
clean interfaces so that each group's sub-part can be independently
evaluated.

* Extend one of the existing labs to add some interesting additional
  functionality to it. Here are a few examples:

  * Lab 2 assessed performance and isolation for the Welp
    application. Investigate the other applications that ship as part
    of DeathStarBench. You can find them all in the DeathStarBench
    subdirectory of the cse453 repository. You can investigate each
    application's performance properties, as well as the performance
    of mixes of these applications. What do you find?  Do mixes of
    applications consolidate better than the same type of application?
    How sensitive is each application to the virtualization technology
    used? How well is performance isolated among mixes of
    applications?

* Design and prototype a lab that could be used in a future iteration
  of this course. Some of the projects listed below are possible
  candidates. For a lab to be feasible at the scale of an entire class
  (versus a single group), it is important to keep the project limited
  and well-contained. You should also describe what you believe should
  be provided as infrastructure versus what is part of the assignment.

* Some of the papers that we have read this quarter propose and
  evaluate an algorithm. For example, Shenango proposes that we can
  take arriving tasks and parcel them out to different cores, stealing
  work to balance load. Equally, in lecture, I have sketched various
  other algorithms to address different issues (e.g., MSI cache
  coherence, or count-min sketch for keeping track of the sizes of the
  largest k connections). Implement a version to see how well it
  works. For Shenango, you could use this to study the impact of
  locality on scheduling - what if some tasks run faster (because of
  better cache behavior) when scheduled onto the same core as some
  previous thread?

* Write a program to determine the hardware performance
  characteristics of the server you are using, such as the latency of
  cache coherence operations between cores, the cost of a TLB miss,
  the extra cost of NUMA memory accesses, etc. Linux has a method to
  limit a particular thread to run on a specific core, and the x86 has
  a way to access a per-processor cycle counter. Given those building
  blocks, and a knowledge of how MSI works, we should be able to test
  how long it takes for a read to data that has been recently modified
  by a different processor, depending on where that processor is
  running. We should be able to compare that to the latency of write
  operations to memory that has been recently read on a different
  core, e.g., using a memory synchronization fence. Skipping the
  fence, we should also be able to determine how large the write
  buffer is - how many writes are needed before the CPU stalls waiting
  for the remote access. Similarly, by skipping around in memory, it
  is possible to drive a TLB into a specific pattern of cache
  misses. Does this cost change when we are running inside a (nested)
  virtual machine? Can you tell if the OS or hypervisor is using
  superpages?

* Write a parallel program for some useful application or library, and
  then measure its scalability (or lack of scalability) as you
  increase the number of processors and/or scale beyond a single
  socket or a single server. As an example, a cuckoo hash table
  (cf. wikipedia) is a good data structure for providing low tail
  latency for reads (relative to chaining) because all reads complete
  with at most two lookups. Hash collisions are resolved on writes to
  preserve low tail latency reads. Are there design alternatives that
  might improve its concurrency?

* Borg has a bin-packing problem that is much like the video game
  tetris. Arriving applications request a certain number of
  processors, amount of memory, rate of disk and network I/O. In a
  real system, resources can be adjusted at runtime, but that can
  require the system to relocate the application to a different server
  if the resources aren't available on the current server, and so we
  can defer that for now. Borg's problem is then to find an assignment
  of applications to servers that satisfies the workload with as
  little wasted resources as possible. For example, that would allow
  the system to turn off some servers, saving energy. What kind of
  assignment algorithm would you implement? More complex is to take
  into account network locality - that applications that do a lot of
  internal communication should be assigned servers in the same
  rack. Applications that have bursty or heavy-tailed workloads should
  be assigned to balance load across racks, so that the energy draw
  (and cooling) per rack is more even. Either way, you will need to
  simulate a diversity of arriving applications to schedule -
  obviously, if all the applications are exactly the same the problem
  is trivial.

* Build a virtual block manager for a zoned SSD. (Equivalently, build
  the flash controller logic that maps virtual blocks to physical
  blocks.) When zones are filled, reclaim space to write new blocks by
  compacting and erasing old zones. To make this more complex, add in
  the constraint of wear levelling - that you are trying to minimize
  the maximum number of times any particular zone is erased and
  re-written. (Once a disk has been erased too many times, it becomes
  unusable, meaning your disk size effectively shrinks, behavior that
  can be non-intuitive for applications and users.) To make this more
  ambitious, build an object store for a zoned SSD, so that the
  objects stored on the SSD could be variable sized. For a workload,
  Figure 10 of the Twitter OSDI 2020 paper has a recent measurement of
  the distribution of object sizes for a cloud workload.

* Write some programs that test how well the system isolates
  performance between different processes, using containers or Linux
  cgroups for isolation. For example, suppose one application
  references a lot of different memory pages, while the other is
  memory intensive, but on a smaller and more constrained set of
  pages. Do they interfere more with each other than with a different
  copy of themselves? Another example are two applications that write
  data to the file system in different ways - e.g., one loops writing
  and then sync'ing a small amount of data, while the other writes a
  large amount of data in a batch before sync'ing. Another possible
  example of performance interference would be an application that
  creates a large number of TCP connections versus another that
  creates only a single TCP connection but sends a large amount of
  data on it, eg., where both are communicating with the same physical
  server. Do they share the network resource fairly?

* Explore the consequence of burstiness for provisioning. As we
  discussed in class, tasks that arrive according to an exponential or
  Poisson process will have better tail latency than those that arrive
  according to a bursty or self-similar arrival process. Devise an
  experiment (or set of experiments) to quantify this effect. For
  example, suppose we want tail latency to be no worse than a factor
  of 10x the best case behavior, 99.9% of the time. What is the limit
  of how much load we can put on the system with Poisson or
  self-similar workloads? Does it matter how bursty the workload is?
  What if the service time (the amount of work to do) is also Poisson
  or self-similar?

* Build an algorithm for wiring a data center using a Clos topology,
  given as inputs the number of servers and switch degree (number of
  inputs and outputs). Use that as input to a packet level simulation
  that can show that your topology is capable of routing a packet
  between any two pairs of servers. How many paths exist between any
  pair of servers? One can extend this to support wire bundling, as
  described in the Jupiter paper, or to allow switches to have
  different link capacities at the edge versus the core. Today, a
  switch might support 32 400Gbps links, or be configured as a top of
  rack switch with 16 400Gbps links up into the network and 64 100Gbps
  links down to the servers within the rack. Another extension would
  be to allow for generating a wiring diagram for a certain amount of
  oversubscription, such as 2x between the top of rack switch and the
  aggregation switch.
