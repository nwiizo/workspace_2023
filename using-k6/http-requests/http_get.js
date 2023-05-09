import http from 'k6/http';

export const options = {
  // Key configurations for avg load test in this section
  stages: [
    { duration: '1m', target: 1 }, // traffic ramp-up from 1 to 10 users over 5 minutes.
  ],
};

export default function() {
  http.get('http://localhost:8080/');
}
