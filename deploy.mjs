import { Client } from 'ssh2';
const conn = new Client();
conn.on('ready', () => {
  conn.exec('cd /opt/apps/pragati && docker compose build web && docker compose up -d web', (err, stream) => {
    let o = '';
    stream.on('data', d => o += d.toString());
    stream.on('close', (code) => {
      console.log('web build (exit ' + code + '):\n' + o);
      conn.end();
    });
  });
}).connect({ host: '192.168.1.2', username: 'hoysala', password: 'MDRSKogunde@32', readyTimeout: 10000 });
