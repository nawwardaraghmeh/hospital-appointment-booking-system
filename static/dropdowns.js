// loadCities fetches the initial list of cities from the REST API
async function loadCities() {
  let r = await fetch('/api/cities');
  let d = await r.json();
  city.innerHTML = d.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
  loadHospitals();
}

// loadHospitals filters the hospital list based on the selected city
async function loadHospitals() {
  let r = await fetch('/api/hospitals?city=' + city.value);
  let d = await r.json();
  hospital.innerHTML = d.map(h => `<option value="${h.id}">${h.name}</option>`).join('');
  loadDepartments();
}

// loadDepartments fetches departments for the chosen hospital 
async function loadDepartments() {
  if (!department) return;
  let r = await fetch('/api/departments?hospital=' + hospital.value);
  let d = await r.json();
  department.innerHTML = d.map(dep => `<option value="${dep.id}">${dep.name}</option>`).join('');
}

// event listeners ensure that when a user changes a selection, the lower levels update automatically.
city.onchange = loadHospitals;
hospital.onchange = loadDepartments;

// initialize the chain when the page loads
loadCities();
