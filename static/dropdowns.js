async function loadCities() {
  let r = await fetch('/api/cities');
  let d = await r.json();
  city.innerHTML = d.map(c => `<option value="${c.id}">${c.name}</option>`).join('');
  loadHospitals();
}

async function loadHospitals() {
  let r = await fetch('/api/hospitals?city=' + city.value);
  let d = await r.json();
  hospital.innerHTML = d.map(h => `<option value="${h.id}">${h.name}</option>`).join('');
  loadDepartments();
}

async function loadDepartments() {
  if (!department) return;
  let r = await fetch('/api/departments?hospital=' + hospital.value);
  let d = await r.json();
  department.innerHTML = d.map(dep => `<option value="${dep.id}">${dep.name}</option>`).join('');
}

city.onchange = loadHospitals;
hospital.onchange = loadDepartments;

loadCities();
