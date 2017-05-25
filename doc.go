//Linda is a background manager to poll jobs from broker and dispatch them to multi workers.
//
//Linda Broker provides a unified API across different broker (queue) services.
//
//Brokers allow you to defer the processing of a time consuming task.
//
//Use ReleaseWithDelay func, you can implement a cron job service.
package linda
