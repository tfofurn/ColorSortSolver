# Color Sort Solver

Ever see those color sorting mobile games?  The starting state is a bunch of tubes, some of which are empty, others are filled with layers of colored liquid or balls or something.  You're only allowed to transfer into a tube where the top color matches the color of what you're pouring into it.  The goal is to pour from tube to tube until each tube is either empty or has only one color in it.  

I bought one of these apps and breezed through most of the levels.  I got stuck for a really long time on one, but eventually solved it.  Just a few levels later, I got really stuck.  That particular app allows buying an additional empty tube, which significantly reduces the difficulty, but what fun is that?  Enter the Color Sort Solver!

# Input format

The program reads CSV files provided as command-line arguments.  Within each file, the top row is interpreted as the names of the tubes, and the rest are understood to be the colors in each tube from top to bottom.  Use any color names you like!  
# Known Issues

**Fixed tube height.**  There isn't any checking of the input file's line count.  If you have anything other than four rows of colors, expect bad behavior.

**No intelligence.**  Brute force FTW!  This program considers nearly every possible pairing of tubes at each stage and does a depth-first search of the solution space.  The only consideration of what's under the top color in a tube is to ascertain whether the tube is already monochromatic.  There's probably a good algorithm for computing fewest-steps solutions with less work.

**Not exhaustive.** There was a bug in the original termination logic that allowed the program to keep running after identifying a solution.  It claimed to have found *billions* of solutions for the inciting incident puzzle.  Solution 102,505,948 was the first to require only 38 steps!  While stumbling over a fast solution is fun when it happens, I'm content to not run my laptop's fan that hard.



# Disclaimers

This is my second attempt at programming in Go and my first success with it.  My code is almost certainly not idiomatic, and probably shouldn't be used as an example.