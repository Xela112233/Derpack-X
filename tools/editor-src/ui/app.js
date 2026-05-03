// Derpack Editor frontend
// Vanilla JS, no framework. Single page with a mod list and modals.

const $ = (sel) => document.querySelector(sel);
const $$ = (sel) => document.querySelectorAll(sel);

// ----- API ---------------------------------------------------------------

async function apiGet(path) {
  const res = await fetch(path);
  if (!res.ok) {
    let detail = '';
    try { detail = (await res.json()).error || ''; } catch {}
    throw new Error(`${res.status}: ${detail || res.statusText}`);
  }
  return res.json();
}

async function apiPost(path, body) {
  const res = await fetch(path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body || {}),
  });
  const data = await res.json().catch(() => ({}));
  if (!res.ok) {
    throw new Error(data.error || `${res.status}: ${res.statusText}`);
  }
  return data;
}

// ----- State -------------------------------------------------------------

let allMods = [];
let filterText = '';

// ----- Status log --------------------------------------------------------

function logStatus(level, message) {
  const log = $('#status-log');
  const li = document.createElement('li');
  const ts = new Date().toLocaleTimeString();
  li.innerHTML = `<span class="ts">${ts}</span><span class="${level}">${escapeHtml(message)}</span>`;
  log.prepend(li);
  // Keep at most 50 entries.
  while (log.children.length > 50) {
    log.removeChild(log.lastChild);
  }
}

function escapeHtml(s) {
  return String(s).replace(/[&<>"']/g, (c) => ({
    '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;',
  }[c]));
}

// ----- Pack info + mod list rendering ------------------------------------

async function loadPack() {
  try {
    const pack = await apiGet('/api/pack');
    $('#pack-name').textContent = pack.name || 'Unknown pack';
    $('#pack-version').textContent = pack.version ? `v${pack.version}` : '';
  } catch (err) {
    logStatus('err', `Failed to load pack info: ${err.message}`);
  }
}

async function loadMods() {
  try {
    allMods = await apiGet('/api/mods');
    renderMods();
  } catch (err) {
    logStatus('err', `Failed to load mods: ${err.message}`);
    $('#mod-list-body').innerHTML = `<tr><td colspan="6" class="empty">Error: ${escapeHtml(err.message)}</td></tr>`;
  }
}

function renderMods() {
  const body = $('#mod-list-body');
  const filtered = filterText
    ? allMods.filter(m =>
        m.slug.toLowerCase().includes(filterText) ||
        (m.name || '').toLowerCase().includes(filterText))
    : allMods;

  $('#mod-count').textContent = `${allMods.length} mods${filterText ? ` (${filtered.length} shown)` : ''}`;

  if (filtered.length === 0) {
    body.innerHTML = `<tr><td colspan="6" class="empty">${
      allMods.length === 0 ? 'No mods yet. Click "+ Add mod" to start.' : 'No mods match the filter.'
    }</td></tr>`;
    return;
  }

  body.innerHTML = filtered.map(modRow).join('');
  // Wire up per-row buttons.
  $$('button[data-action]').forEach(btn => {
    btn.addEventListener('click', () => {
      const action = btn.dataset.action;
      const slug = btn.dataset.slug;
      handleRowAction(action, slug);
    });
  });
}

function modRow(m) {
  const sourceTag = m.source
    ? `<span class="source-tag ${m.source}">${m.source}</span>`
    : '';
  const pinIcon = m.pinned
    ? '<span class="pin-icon pinned" title="Pinned">📌</span>'
    : '<span class="pin-icon unpinned" title="Not pinned">📌</span>';

  return `
    <tr>
      <td><code>${escapeHtml(m.slug)}</code></td>
      <td>${escapeHtml(m.name || '')}</td>
      <td>${escapeHtml(m.side || 'both')}</td>
      <td>${sourceTag}</td>
      <td class="col-pinned">${pinIcon}</td>
      <td class="col-actions">
        <div class="row-actions">
          <button data-action="${m.pinned ? 'unpin' : 'pin'}" data-slug="${escapeHtml(m.slug)}">
            ${m.pinned ? 'Unpin' : 'Pin'}
          </button>
          <button data-action="remove" data-slug="${escapeHtml(m.slug)}">Remove</button>
        </div>
      </td>
    </tr>
  `;
}

// ----- Row actions -------------------------------------------------------

async function handleRowAction(action, slug) {
  switch (action) {
    case 'pin':
      await doPin(slug, true);
      break;
    case 'unpin':
      await doPin(slug, false);
      break;
    case 'remove':
      askConfirm(
        `Remove ${slug}?`,
        `This will delete mods/${slug}.pw.toml and re-index.`,
        () => doRemove(slug),
      );
      break;
  }
}

async function doPin(slug, pinned) {
  try {
    const r = await apiPost('/api/mods/pin', { slug, pinned });
    if (r.ok) {
      logStatus('ok', `${pinned ? 'Pinned' : 'Unpinned'} ${slug}`);
      await loadMods();
    } else {
      logStatus('err', `${pinned ? 'Pin' : 'Unpin'} failed for ${slug}: ${r.error || 'unknown error'}`);
    }
  } catch (err) {
    logStatus('err', err.message);
  }
}

async function doRemove(slug) {
  try {
    const r = await apiPost('/api/mods/remove', { slug });
    if (r.ok) {
      logStatus('ok', `Removed ${slug}`);
      await loadMods();
    } else {
      logStatus('err', `Remove failed for ${slug}: ${r.error || 'unknown error'}`);
    }
  } catch (err) {
    logStatus('err', err.message);
  }
}

// ----- Add mod modal -----------------------------------------------------

function openAddModal() {
  $('#add-source').value = 'mr';
  $('#add-slug').value = '';
  $('#add-side').value = 'both';
  $('#add-modal').classList.remove('hidden');
  $('#add-slug').focus();
}

function closeAddModal() {
  $('#add-modal').classList.add('hidden');
}

async function submitAddMod() {
  const source = $('#add-source').value;
  const slug = $('#add-slug').value.trim();
  const side = $('#add-side').value;
  if (!slug) {
    return;
  }

  $('#add-confirm').disabled = true;
  $('#add-confirm').textContent = 'Adding...';
  try {
    const r = await apiPost('/api/mods/add', { source, slug, side });
    if (r.ok) {
      logStatus('ok', `Added ${slug} (${source})`);
      closeAddModal();
      await loadMods();
    } else {
      logStatus('err', `Add failed for ${slug}: ${r.error || 'unknown error'}`);
    }
  } catch (err) {
    logStatus('err', err.message);
  } finally {
    $('#add-confirm').disabled = false;
    $('#add-confirm').textContent = 'Add';
  }
}

// ----- Confirm modal -----------------------------------------------------

let confirmCallback = null;

function askConfirm(title, message, onConfirm) {
  $('#confirm-title').textContent = title;
  $('#confirm-message').textContent = message;
  $('#confirm-modal').classList.remove('hidden');
  confirmCallback = onConfirm;
}

function closeConfirm() {
  $('#confirm-modal').classList.add('hidden');
  confirmCallback = null;
}

// ----- Wire up -----------------------------------------------------------

document.addEventListener('DOMContentLoaded', () => {
  // Top bar
  $('#refresh-btn').addEventListener('click', () => {
    loadPack();
    loadMods();
    logStatus('ok', 'Refreshed');
  });

  // Search
  $('#search').addEventListener('input', (e) => {
    filterText = e.target.value.toLowerCase();
    renderMods();
  });

  // Add modal
  $('#add-mod-btn').addEventListener('click', openAddModal);
  $('#add-cancel').addEventListener('click', closeAddModal);
  $('#add-confirm').addEventListener('click', submitAddMod);
  $('#add-slug').addEventListener('keydown', (e) => {
    if (e.key === 'Enter') submitAddMod();
  });

  // Confirm modal
  $('#confirm-cancel').addEventListener('click', closeConfirm);
  $('#confirm-ok').addEventListener('click', () => {
    const cb = confirmCallback;
    closeConfirm();
    if (cb) cb();
  });

  // Close modals when clicking outside
  $$('.modal').forEach(modal => {
    modal.addEventListener('click', (e) => {
      if (e.target === modal) modal.classList.add('hidden');
    });
  });

  // Initial load
  loadPack();
  loadMods();
});
