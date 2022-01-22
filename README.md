# Color Sort Solver

Ever see those color sorting mobile games?  The starting state is a bunch of tubes, some of which are empty, others are filled with layers of colored liquid or balls or something.  You're only allowed to transfer into a tube where the top color matches the color of what you're pouring into it.  The goal is to pour from tube to tube until each tube is either empty or is filled with a single color.  

I bought one of these apps and breezed through most of the levels.  I got stuck for a really long time on one, but eventually solved it.  Just a few levels later, I got really stuck.  That particular app allows buying an additional empty tube, which significantly reduces the difficulty, but what fun is that?  Enter the Color Sort Solver!

# Usage

If you want a solution for a particular puzzle, specify that file on the command line:

> `colorsortsolver sample/inciting-incident.csv`

If no arguments are provided, the program will automatically attempt every file in the `samples` folder.  It will print timings and solution counts, but not the solutions.  This is how I assess whether code changes impact solving performance.  

# Input format

The program reads CSV files.  Within each file, the top row is interpreted as the names of the tubes, and the rest of the rows are understood to be the colors in each tube from top to bottom.  Use any tube and color names you like!  Be sure to provide tube names at the tops of empty columns to reflect the empty tubes at the start!  

The `sample` directory contains several puzzles, mostly those I wasn't patient enough to solve as a human.  The app on my phone displays two rows of tubes, so I started with a `R[row]T[tube]` scheme.  I later switched to `Top[tube]` and `Bot[tube]`.  

# Known Issues

**Fixed tube height.**  There isn't any checking of the input file's line count.  If you have anything other than four rows of colors, expect bad behavior.

**No intelligence.**  Brute force for the win!  This program considers nearly every possible pairing of tubes at each stage and does a depth-first search of the solution space.  The only consideration of what's under the top color in a tube is to ascertain whether the tube is already monochromatic.  There's probably a good algorithm for computing fewest-steps solutions with less work.

**Not exhaustive.** There was a bug in the original termination logic that allowed the program to keep running after identifying a solution.  It claimed to have found *billions* of solutions for the inciting incident puzzle.  Solution 102,505,948 was the first to require only 38 steps!  While stumbling over a fast solution is fun when it happens, I'm content to not run my laptop's fan that hard.

**Unnecessary constraints.** The solver will only consider steps that pour out all contiguous layers of a color.  Very occasionally as a human, I've split a multi-layer segment into two different tubes.  So far, the samples I've run through the program haven't needed this technique to find a solution, so it's not considered.

**Disappointing multithreading.** The program splits the workload at the beginning, assigning each worker a partial solution to use as the basis for the solve.  It's frequently the case that most of the solvers will report a solution while one keeps laboring away.  I would prefer to keep all workers active to the end, but haven't figured out how to adequately communicate amongst the workers when the last partial solution has been queued.

# Disclaimers

This is my second attempt at programming in Go and my first success with it.  My code is almost certainly not idiomatic, and probably shouldn't be used as an example.
