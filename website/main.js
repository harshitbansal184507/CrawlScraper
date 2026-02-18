const observer = new IntersectionObserver((entries) => {
  entries.forEach(e => {
    if (e.isIntersecting) {
      e.target.classList.add('visible');
    }
  });
}, { threshold: 0.1 });

document.querySelectorAll('.reveal').forEach(el => observer.observe(el));

//switching bw single and mutlple 
function switchTab(tab) {
  document.getElementById('tab-single').classList.add('hidden');
  document.getElementById('tab-multi').classList.add('hidden');
  document.getElementById('tab-' + tab).classList.remove('hidden');

  document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
  event.target.classList.add('active');
}