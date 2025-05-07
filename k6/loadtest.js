import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 10,
  duration: '10s',
};

export default function () {
  const loginRes = http.post('http://server:8080/login', JSON.stringify({
    username: 'testuser',
    password: 'testpass',
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  const token = loginRes.json('token');

  check(loginRes, {
    'login succeeded': (res) => res.status === 200 && token,
  });

  const privateRes = http.get('http://server:8080/private', {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });

  check(privateRes, {
    'access to private endpoint succeeded': (res) => res.status === 200,
  });

  sleep(1);
}