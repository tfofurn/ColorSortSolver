# Color Sort Solver

Ever see those color sorting mobile games?  The starting state is a bunch of tubes, some of which are empty, others are filled with layers of colored liquid or balls or something.  You're only allowed to transfer into a tube where the top color matches the color of what you're pouring into it.  The goal is to pour from tube to tube until each tube is either empty or has only one color in it.  

I bought one of these apps and breezed through most of the levels.  I got stuck for a really long time on one, but eventually solved it.  Just a few levels later, I got really stuck.  That particular app allows buying an additional empty tube, which significantly reduces the difficulty, but what fun is that?  Enter the Color Sort Solver.

# Known Issues

**Hardcoded color names.** Makes data entry a pain.

**Way too much output.**  In addition to producing multiple trivially-different versions of each solution, there are still a bunch of debugging print statements lingering.

**No intelligence.**  Brute force FTW!  This program considers literally every possible pairing of tubes at each stage and does a depth-first search of the solution space.  The only consideration of what's under the top color in a tube is to ascertain whether the tube is already monochromatic.  There's probably a good algorithm for computing fewest-steps solutions with less work.

**Runs forever.** While the first solution of the inciting puzzle pops out within seconds, it takes a while to exhaustively search the solution space.  I think it's still emitting unique solutions all the way through; I originally thought I designed it to bail after outputting one. :-)  I also haven't implemented signalling for when the worker threads are finished.



# Disclaimers

This is my second attempt at programming in Go and my first success with it.  My code is almost certainly not idiomatic, and probably shouldn't be used as an example.