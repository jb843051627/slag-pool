async function fetchJSON(url) {
    try {
        const resp = await fetch(url);
        if (!resp.ok) return null;
        return await resp.json();
    } catch (e) {
        console.error('fetch error:', url, e);
        return null;
    }
}

async function loadPools() {
    const pools = await fetchJSON('/api/pools');
    if (!pools || pools.length === 0) {
        document.getElementById('pool-list').innerHTML = '<div class="pool-card">暂无渣池数据</div>';
        return;
    }
    const container = document.getElementById('pool-list');
    container.innerHTML = '';
    for (const pool of pools) {
        const status = await fetchJSON(`/api/pools/${pool.id}/status`);
        const card = document.createElement('div');
        card.className = 'pool-card';
        const temp = status ? status.temp.toFixed(1) : '--';
        const flow = status ? status.flow.toFixed(1) : '--';
        const alerts = status ? status.alert_count : 0;
        card.innerHTML = `
            <div class="name">${pool.name}</div>
            <div class="meta">${pool.location} | <span class="status-${pool.status}">${pool.status}</span></div>
            <div class="metrics">
                <div class="metric"><span class="label">温度</span> <span class="value">${temp}°C</span></div>
                <div class="metric"><span class="label">流量</span> <span class="value">${flow}L/min</span></div>
                <div class="metric"><span class="label">告警</span> <span class="value">${alerts}</span></div>
            </div>`;
        container.appendChild(card);
    }
}

async function loadAlerts() {
    const alerts = await fetchJSON('/api/alerts?status=active');
    const tbody = document.querySelector('#alert-table tbody');
    if (!alerts || alerts.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5">暂无告警</td></tr>';
        return;
    }
    tbody.innerHTML = alerts.map(a => `
        <tr class="alert-${a.level}">
            <td>${a.pool_id}</td><td>${a.type}</td><td>${a.level}</td>
            <td>${a.message}</td><td>${new Date(a.created_at).toLocaleString()}</td>
        </tr>`).join('');
}

async function loadBatches() {
    const pools = await fetchJSON('/api/pools');
    const tbody = document.querySelector('#batch-table tbody');
    if (!pools || pools.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6">暂无批次</td></tr>';
        return;
    }
    let rows = [];
    for (const pool of pools) {
        const batches = await fetchJSON(`/api/pools/${pool.id}/batches`);
        if (batches && batches.length > 0) {
            rows.push(...batches.map(b => `
                <tr>
                    <td>${b.batch_number || b.id}</td><td>${pool.name}</td>
                    <td class="status-${b.status}">${b.status}</td>
                    <td>${b.target_temp}°C</td><td>${b.actual_temp}°C</td>
                    <td>${new Date(b.start_time).toLocaleString()}</td>
                </tr>`));
        }
    }
    tbody.innerHTML = rows.length > 0 ? rows.join('') : '<tr><td colspan="6">暂无批次</td></tr>';
}

async function loadReadings() {
    const pools = await fetchJSON('/api/pools');
    const tbody = document.querySelector('#reading-table tbody');
    if (!pools || pools.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5">暂无读数</td></tr>';
        return;
    }
    let rows = [];
    const pool = pools[0];
    const readings = await fetchJSON(`/api/pools/${pool.id}/latest`);
    if (readings && readings.length > 0) {
        rows = readings.map(r => `
            <tr>
                <td>${r.sensor_id}</td><td>${r.value.toFixed(2)}</td>
                <td>${r.unit}</td><td>${r.quality}</td>
                <td>${new Date(r.timestamp).toLocaleString()}</td>
            </tr>`);
    }
    tbody.innerHTML = rows.length > 0 ? rows.join('') : '<tr><td colspan="5">暂无读数</td></tr>';
}

async function refresh() {
    await Promise.all([loadPools(), loadAlerts(), loadBatches(), loadReadings()]);
}

refresh();
setInterval(refresh, 5000);
