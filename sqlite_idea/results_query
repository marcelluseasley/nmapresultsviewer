SELECT hostdata.uuid,
hostdata.ip as ip,
hostdata.host_state as host_state,
hostdata.reason as h_reason,
hostdata.hostname as hostname,
portdata.port as port,
portdata.state as p_state,
portdata.reason as p_reason,
portdata.service as service,
portdata.method as method
FROM hostdata
INNER JOIN portdata 
ON hostdata.uuid = portdata.uuid 
where hostdata.uuid = '1af6effd-ff62-4130-8058-ebe5f218bb27'
AND hostdata.ip = portdata.ip;