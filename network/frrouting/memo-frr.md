# Memo: FRRouting

## dataplane context objects

"dataplane context objects" is something you should understand when first working on FRRouting dataplane.

> [FRR latest documentation >> Zebra >> Design](https://docs.frrouting.org/projects/dev-guide/en/latest/zebra.html#design)
> With our dataplane abstraction, we create a queue of dataplane context objects for the messages we want to send to the kernel. In a separate pthread, we loop over this queue and send the context objects to the appropriate dataplane. A batching enhancement tightly integrates with the dataplane context objects so they are able to be batch sent to dataplanes that support it.

