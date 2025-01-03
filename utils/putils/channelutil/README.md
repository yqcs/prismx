# channelutil

Channel utils provides a set of utilities for mux and demux channels specifically `cloning` and `joining` channels.
This is useful when you want to send the same data to multiple channels or when you want to send data from multiple channels to a single channel.

### Cloning

A Simple approach for cloning a channel is to send data received from a channel to multiple channels by looping over the channels and sending the data to each channel.
This is a blocking operation and totally depends on what/How data of the channel is being consumed. if anyone of the channel is blocked for some reason then the other channels will also be blocked.
and this happens often if the number of channels is large.

To overcome this `CloneChannels` implements a relay channels which receives data from a single channel and sends it to multiple channels (5 by default) in a non-blocking way using select statement.
since select statement picks the channel that is ready to receive data, it is non-blocking and if data is sent to x channel data is put into buffer of other channels and auto drain is triggered
when a buffer of particular channel is full.

```	
		If sinks > 5
		relay channels are used that relay data from root node to leaf node (i.e in this case channel)

		1. sinks are grouped into 5 with 1 relay channel for each group
		2. Each group is passed to worker
		3. Relay are fed to Clone i.e Recursion
	
	
			Ex:
                                    $ 			 <-  Source Channel
	                          /   \
	                         $     $			 <-  Relay Channels
			        / \    / \
			       $   $  $   $		 <-  Leaf Channels (i.e Sinks)

		*Simplicity purpose 2 childs are shown for each nodebut each node(except root node) has 5 childs
```	

### Joining

Joining is the opposite of cloning. It receives data from multiple channels and sends it to a single channel and again simple approach is to loop over the channels and send data to the single channel
but this is blocking if the single channel is blocked for some reason then all the channels will be blocked and this happens often if the number of channels is large. If channels are known in advance then
select statement will be better and easier approach. But if number of channels are not known in advance or it is dynamic then `JoinChannels` is the way to go.

Go Standard library provides [reflect.SelectDir](https://pkg.go.dev/reflect#SelectDir) for this purpose but it is known to be very slow and inefficient. `JoinChannels` implements a relay channel which receives data from multiple channels and sends it to a single channel in a non-blocking way using select statement when a source channel is completely drained it is set to nil and select statement will not pick that case again

```
	
		If sources > 5
		relay channels are used that relay data from leaf nodes to root node (i.e in this case channel)

		1. sources are grouped into 5 with 1 relay channel for each group
		2. Each group is passed to worker
		3. Relay are fed to Join i.e Recursion
	
	
		Ex:
			$   $ $   $		 <-  Leaf Channels (i.e Sources)
			 \ /   \ /
		          $  	$		 <-  Relay Channels
			   \   /
			     $           <- Sink Channel

		*Simplicity purpose 2 childs are shown for each node but each node has 5 childs
	
```

### References

- https://medium.com/justforfunc
