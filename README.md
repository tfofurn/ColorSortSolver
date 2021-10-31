# Color Sort Solver

Ever see those color sorting mobile games?  The starting state is a bunch of tubes, some of which are empty, others are filled with layers of colored liquid or balls or something.  You're only allowed to transfer into a tube where the top color matches the color of what you're pouring into it.  The goal is to pour from tube to tube until each tube is either empty or has only one color in it.  

I bought one of these apps and breezed through most of the levels.  I got stuck for a really long time on one, but eventually solved it.  Just a few levels later, I got really stuck.  That particular app allows buying an additional empty tube, which significantly reduces the difficulty, but what fun is that?  Enter the Color Sort Solver.

# Known Issues

**Hardcoded color names and file path.** Makes data entry a pain.

**No intelligence.**  Brute force FTW!  This program considers literally every possible pairing of tubes at each stage and does a depth-first search of the solution space.  The only consideration of what's under the top color in a tube is to ascertain whether the tube is already monochromatic.  There's probably a good algorithm for computing fewest-steps solutions with less work.

**Not exhaustive.** There was a bug in the original termination logic that allowed the program to keep running after identifying a solution.  It claimed to have found *billions* of solutions for the inciting incident puzzle.  Solution 102,505,948 was the first to require only 38 steps!  While stumbling over a fast solution is fun when it happens, I'm content to not run my laptop's fan that hard.



# Disclaimers

This is my second attempt at programming in Go and my first success with it.  My code is almost certainly not idiomatic, and probably shouldn't be used as an example.