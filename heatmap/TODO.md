# Refactor getter functions
IDEA:
    Rather than sending raycasts and positions individually, nest this information inside of the player structs and 
    send the player structs to the app in svelte

    Problems: 
    The original player definitions contain in-game coordinates and not canvas coordinates, would need to be re-defined
    to canvas coordinates which would destroy all code

# Rewrite all vector code to use in-game info
IDEA:
    Current code uses canvas data which makes translating between svelte and js annoying,
    it can be rewritten so that all calculations are done with in-game data, and is translated to canvas-readable data
    at the very end.

    Will make code more readable and extendable.

    Problems:
    All code needs to be rewritten and thought-through.

# Make frontend more presentable and easy-to-use

# Refactor all frontend code

# Write unit tests
IDEA:
    I often have to test each section of functions individually etc. It would be better to have this done automatically
    each time so I can diagnose a problem much faster

    As well, I will be able to tell when a change screws something up elsewhere rather than waiting for the problem to creep in
    eventually.

    Problems:
    Not all sections of each function may be able to be tested

    Maybe need to refactor functions completely to have effective unit tests