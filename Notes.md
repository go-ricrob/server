# The Solver

This document shall collect some ideas and concepts regarding the solver

## Rules of the game

The board consists of 4 tiles each with 2 faces providing a total of 96 different combinations.
On the board, 4 or 5 robots of unique color are placed that can move in horizontal or vertical direction.
However, the robots can't stop before they reach the next obstacle - either a wall or another robot.

There are 4 targets for each of the four colors red, green, blue and yellow.
These targets have to be occupied by the robot with the respective color.

There is one further target named cosmic whirl or cosmic card. 
This target can be occupied by any robot, the color doesn't matter. 
This applies especially if the game is played with an additional silver robot.

Before entering the target field, the respective robot has to ricochet at least once.
To ricochet means that two consecutive moves form a rectangular path.

Some people use an additional rule that rebound is not allowed, i.e. moving a robot back and forth.
However, according to the standard rules, this is not forbidden and might be reasonable if one or more robots of different color have moved in between. In general, considering the same situation more than once doesn't make any sense.

## Problem complexity

the puzzle "Ricochet Robots" has been proven to be NP hard, see [References](https://github.com/stfnmllr/go-ricrob/blob/main/Notes.md#references)

What does that mean? There does not seem to be a sophisticated way of solving each possible puzzle with a minimum number of robot moves - it's no shame trying every possible move in a brute force way

Here, a breadth first tree traversal approach is chosen, i.e. the nodes within the tree are examined layer by layer. When a solution is found, it's safe to say that there is no solution available with less robot moves.


## Data Structures

### Node

A node instance represents a specific situation, i.e. the positions of all robots that take part in a game

It's necessary to check whether a situation has already been considered within a sequence of moves. 
The parent node instance can be reached via pointer

As there might be many node instances in a puzzle, the memory footprint of the node structure shall be as small as possible

The parent pointer takes 8 Bytes. Using a different memory model with e.g. only 4 Byte pointers would limit the addressable memory to 4 GiB which is not an option.

The robot positions are stored in a packed byte array, each robot position occupies one Byte
On a 64 Bit architecture, an 8 Byte wide integer can be used without penalty

It might be possible to use an index of 4 Bytes width to indirectly address the robot position payload. 
This approach could make use of the memory without any losses
However this seems to require a contigous, linearly addressable memory area with 16 .. 20 GiB in size


### Queue

A queue is used for temporarily storing nodes. It gets obsolete as soon as the respective child node are created.


## The Solving Process / Procedure

1. create the root node from the initial position of all robots, parent pointer gets value nil and put it in the parent queue

2. move each robot in every possible direction
determine the new position of the moving robot from the reach map and the positions of the other robots

2a. if the new position is equal to the current position the robot is blocked, so continue with the next move

2b. if the new position is not equal to the current position check if a solution is found:

* the target field is occupied by a robot

* the color of the robot matches the color of the target field or the target field is the cosmic whirl

* the robot on the target field has ricochet at least once before

If a solution is found, stop the solving process and return the sequence of robot moves 

If a solution is not found create a child node, the parent pointer gets the address of the root node and finally, put the child node in the child queue

As soon as all moves of a node are processed go on with the next node from the parent queue

As soon as all nodes from the parent queue are processed, make the child queue the parent queue and continue on the next depth level

## Iteration vs. Recursion

As the tree of nodes is a recursive data structure, the traversal can be implemented as recursion in a quite natural way. However, every recursive processing can be made an iterative one.

Advantages of an iterative approach seem to be easier handling and less usage of stack space

## References

https://www.sciencedirect.com/science/article/abs/pii/S1571065306000631
Randolphs Robot Game is NP-hard!
Birgit Engels	Zentrum für angewandte Informatik, University of Cologne, Germany
Tom Kamphans	Institut für Informatik I, University of Bonn, Germany

https://boardgamegeek.com/thread/308557/hardest-ricochet-robots-problem

## Open Questions

_Solve needs to be save to be called concurrently is used as plugin function_

Shall there be exactly one Solve instance that is shared by several concurrent contexts or
shall each caller get its own individual Solve instance?

From a functional point of view, it does make sense to share a Solve instance between
several caller contexts considering the same puzzle.

When providing several Solver instances that are independent from each other the resource consumption
might get critical up to the point where no instance is able to finish its task

The resoure consumption especially memory allocation depends on the problem to solve and is a priori unknown.
Therefore, reserving all available system resources for one active solver instance at a time might be a better
approach than limiting resources per solver by e.g. quotas.


