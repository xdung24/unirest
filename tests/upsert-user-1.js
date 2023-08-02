import http from "k6/http";

// Test configuration
export const options = {
  maxRedirects: 1 ,
  thresholds: {
      'http_req_duration{status:0}': ['max>=0'],
      'http_req_duration{status:200}': ['max>=0'],
      'http_req_duration{status:400}': ['max>=0'],
      'http_req_duration{status:500}': ['max>=0'],
      'http_req_duration{status:502}': ['max>=0'],
      'http_req_duration{method:POST}': ['max>=0'],
  },  
  summaryTrendStats: ['min', 'med', 'avg', 'p(90)', 'p(95)', 'max', 'count'],
  discardResponseBodies: true,
  // Ramp the number of virtual users up and down
  stages: [
    { duration: "30s", target: 10 },
    { duration: "30s", target: 20 },
    { duration: "20s", target: 0 },
  ],
};

// Simulated user behavior
export default function () {
  const url = "http://localhost:8000/ns/users/1";
  const payload = JSON.stringify({
      firstName: "jack",
      lastName: "neverdie",
      age: 20  
  });
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiSm9obiBEb2UiLCJqdGkiOiJqb2huZCIsImlhdCI6MTUxNjIzOTAyMn0.wZTch4SZRqI3M2gMUe_t_rUjpxR0XHxecbV7XsVAWHxGBwJe1QqCoDuD6Y0DGHYKx6bAZ8dcEymo9NcXuhMit10VDLWVLptgttXbQ57VBGIXBhaCd5S2emJxu29LaN_yVJCtkTVYRhSc5HcLN6SjkUoEjhr0ezlI_wmbpM7TL7J4Kx87W7oMmLNNrUlbXzaDeHfXmvbNMxfR-x4N25WPnO3ifGIykHk8mUWeRKUprHSiWfhzjx3ZUCCrq1WyxefuJKXx0ONL1NQ8SxGRa3EECuxfuww0Ic3oCJrJR7AMjbOv9o99UZBpqMJk9A2BKxW3rxwAdcr5Wh4kQBwRXkIELg'
    },
  };
  http.post(url, payload, params);
}