import http from "k6/http";

// Test configuration
export const options = {
  thresholds: {
    // Assert that 99% of requests finish within 1000ms.
    http_req_duration: ["p(99) < 1000"],
  },
  // Ramp the number of virtual users up and down
  stages: [
    { duration: "30s", target: 10 },
    { duration: "30s", target: 20 },
    { duration: "20s", target: 0 },
  ],
};

// Simulated user behavior
export default function () {
  let res = http.get("http://localhost:8000/ns/users/1");
}