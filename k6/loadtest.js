import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 20,
  duration: '10s',
};

export default function () {
  const res = http.get('http://server:8080/health');

  check(res, {
    'status is 200': (r) => r.status === 200,
    'body is OK': (r) => r.body.includes('OK'),
  });

  sleep(1);
}