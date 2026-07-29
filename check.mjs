import { Client } from 'ssh2';
const conn = new Client();
conn.on('ready', () => {
  conn.exec('docker logs pragati-api-1 --tail 20 2>&1; echo "---END---"', (err, stream) => {
    let o = '';
    stream.on('data', d => o += d);
    stream.on('close', () => {
      console.log(o);
      // Also check its start time
      conn.exec('docker inspect pragati-api-1 --format "{{.State.StartedAt}}"', (err2, stream2) => {
        let o2 = '';
        stream2.on('data', d => o2 += d);
        stream2.on('close', () => {
          console.log('Started:', o2.trim());
          conn.end();
        });
      });
    });
  });
}).connect({ host:'192.168.1.2', username:'hoysala', password:'MDRSKogunde@32', readyTimeout:10000 });
