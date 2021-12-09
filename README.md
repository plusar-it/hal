# hal
The consumer producer pattern implemented is: one consumer and many producers. The number of producers is configurable in the settings.

The triggers are parsed and loaded in memory by the command package. There is a check for duplicate triggers, that have the same contract address.

The producer package is responsible for:
- start each indvidual go-routine
- map the trigger array in different slices, one for each routine
- trigger processing function calls the Multicaller with batched requests; the batch number is configurable
- publish the even block numbers in the syncronisation channel
- close the syncronisation channel after all the triggers have been processed

The consumer package is responsible for:
- consuming the block numbers from the syncronisation channel
- call the underlying DB layer to save the block number
- consumer exists when the channel is closed ( by the producer)