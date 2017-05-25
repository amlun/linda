package linda

const (
	//SizeScript = `return redis.call('llen', KEYS[1]) + redis.call('zcard', KEYS[2]) + redis.call('zcard', KEYS[3])`

	// PopJobScript
	// -- Pop the first job off of the queue...
	PopJobScript = `local job = redis.call('lpop', KEYS[1])
		if(job ~= false) then
			-- place job on the reserved queue...
			redis.call('zadd', KEYS[2], ARGV[1], job)
		end
		return job`

	// ReleaseScript
	// -- Remove the job from the current queue...
	ReleaseScript = `redis.call('zrem', KEYS[2], ARGV[1])
		-- Add the job onto the "delayed" queue...
		redis.call('zadd', KEYS[1], ARGV[2], ARGV[1])
		return true`

	// MigrateJobsScript
	//-- Get all of the jobs with an expired "score"...
	MigrateJobsScript = `local val = redis.call('zrangebyscore', KEYS[1], '-inf', ARGV[1])
		-- If we have values in the array, we will remove them from the first queue
		-- and add them onto the destination queue in chunks of 100, which moves
		-- all of the appropriate jobs onto the destination queue very safely.
		if(next(val) ~= nil) then
    			redis.call('zremrangebyrank', KEYS[1], 0, #val - 1)
    			for i = 1, #val, 100 do
        			redis.call('rpush', KEYS[2], unpack(val, i, math.min(i+99, #val)))
    			end
		end
		return val`
)
