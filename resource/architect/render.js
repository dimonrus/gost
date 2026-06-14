// ============================================
// STATE MANAGEMENT
// ============================================
let blocks = [];
let connections = [];
let nextId = 1;
let draggedBlock = null;
let dragOffset = {x: 0, y: 0};
let isPanning = false;
let panStart = {x: 0, y: 0};
let resizingBlock = null;
let resizeStart = {x: 0, y: 0, w: 0, h: 0};

// Zoom state
let zoom = 1.0;
const MIN_ZOOM = 0.25;
const MAX_ZOOM = 3.0;
const ZOOM_STEP = 0.1;

// Undo/Redo state
let history = [];
let historyIndex = -1;
const MAX_HISTORY = 20;
let isUndoingRedoing = false;
let saveTimeout = null;

// LocalStorage key
const STORAGE_KEY = 'go_visual_architect_state';

const canvas = document.getElementById('canvas');
const canvasContainer = document.getElementById('canvasContainer');
const emptyState = document.getElementById('emptyState');
const zoomLevelEl = document.getElementById('zoomLevel');

const COLOR_PALETTE = [
    // Blues (10)
    '#58a6ff', '#388bfd', '#0969da', '#1f6feb', '#005cc5',
    '#79c0ff', '#a5d6ff', '#7dd3fc', '#38bdf8', '#0ea5e9',
    // Greens (10)
    '#3fb950', '#2ea043', '#1f883d', '#1a7f37', '#116329',
    '#56d364', '#7ee787', '#4ade80', '#22c55e', '#16a34a',
    // Yellows & Oranges (10)
    '#d29922', '#e3b341', '#9e6a03', '#7d4e00', '#845306',
    '#f59e0b', '#fbbf24', '#fcd34d', '#facc15', '#eab308',
    // Reds & Pinks (10)
    '#f85149', '#ff7b72', '#da3633', '#b42318', '#cf222e',
    '#ef4444', '#f87171', '#fca5a5', '#fb7185', '#f43f5e',
    // Purples & Violets (10)
    '#a371f7', '#bc8cff', '#8957e5', '#6e40c9', '#8250df',
    '#c084fc', '#a78bfa', '#8b5cf6', '#7c3aed', '#6d28d9',
    // Teals & Cyans (10)
    '#2dd4bf', '#14b8a6', '#0d9488', '#0f766e', '#115e59',
    '#5eead4', '#99f6e4', '#06b6d4', '#0891b2', '#0e7490',
    // Warm & Earthy (10)
    '#92400e', '#78350f', '#c2410c', '#ea580c', '#a8a29e',
    '#78716c', '#57534e', '#44403c', '#292524', '#9a3412',
    // Cool Grays & Slate (10)
    '#94a3b8', '#64748b', '#475569', '#334155', '#1e293b',
    '#cbd5e1', '#e2e8f0', '#6b7280', '#4b5563', '#374151'
];

const DEFAULT_COLORS = {
    interface: '#58a6ff',
    struct: '#3fb950',
    method: '#d29922'
};

// ============================================
// UNDO / REDO / LOCALSTORAGE
// ============================================

function saveState() {
    if (isUndoingRedoing) return;

    // Remove any future history if we're not at the end
    if (historyIndex < history.length - 1) {
        history = history.slice(0, historyIndex + 1);
    }

    const state = {
        blocks: JSON.parse(JSON.stringify(blocks)),
        nextId: nextId,
        zoom: zoom,
        timestamp: Date.now()
    };

    history.push(state);

    // Keep only MAX_HISTORY entries
    if (history.length > MAX_HISTORY) {
        history.shift();
    } else {
        historyIndex++;
    }

    updateUndoRedoButtons();
    debouncedLocalStorageSave();
}

function debouncedLocalStorageSave() {
    if (saveTimeout) clearTimeout(saveTimeout);
    saveTimeout = setTimeout(() => {
        saveToLocalStorage();
    }, 500);
}

function saveToLocalStorage() {
    const state = {
        blocks: blocks,
        nextId: nextId,
        zoom: zoom,
        history: history,
        historyIndex: historyIndex
    };
    try {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(state));
    } catch (e) {
        console.warn('Failed to save to localStorage:', e);
    }
}

function loadFromLocalStorage() {
    try {
        const saved = localStorage.getItem(STORAGE_KEY);
        if (saved) {
            const state = JSON.parse(saved);
            if (state.blocks && state.blocks.length > 0) {
                blocks = state.blocks;
                nextId = state.nextId || 1;
                zoom = state.zoom || 1.0;
                if (state.history && state.history.length > 0) {
                    history = state.history;
                    historyIndex = state.historyIndex || history.length - 1;
                }
                return true;
            }
        }
    } catch (e) {
        console.warn('Failed to load from localStorage:', e);
    }
    return false;
}

function restoreState(state) {
    isUndoingRedoing = true;

    // Clear current canvas
    document.querySelectorAll('.block').forEach(el => el.remove());
    document.querySelectorAll('.connection-line').forEach(el => el.remove());

    blocks = JSON.parse(JSON.stringify(state.blocks));
    nextId = state.nextId;

    // Re-render all blocks
    blocks.forEach(block => {
        renderBlock(block);
    });

    if (blocks.length === 0) {
        emptyState.style.display = 'block';
    } else {
        emptyState.style.display = 'none';
    }

    updateConnections();
    updateUndoRedoButtons();

    isUndoingRedoing = false;
}

function undo() {
    if (historyIndex > 0) {
        historyIndex--;
        restoreState(history[historyIndex]);
        saveToLocalStorage();
    }
}

function redo() {
    if (historyIndex < history.length - 1) {
        historyIndex++;
        restoreState(history[historyIndex]);
        saveToLocalStorage();
    }
}

function updateUndoRedoButtons() {
    const undoBtn = document.getElementById('undoBtn');
    const redoBtn = document.getElementById('redoBtn');

    if (undoBtn) {
        undoBtn.disabled = historyIndex <= 0;
        undoBtn.style.opacity = historyIndex <= 0 ? '0.4' : '1';
        undoBtn.style.cursor = historyIndex <= 0 ? 'not-allowed' : 'pointer';
    }
    if (redoBtn) {
        redoBtn.disabled = historyIndex >= history.length - 1;
        redoBtn.style.opacity = historyIndex >= history.length - 1 ? '0.4' : '1';
        redoBtn.style.cursor = historyIndex >= history.length - 1 ? 'not-allowed' : 'pointer';
    }
}

function clearHistory() {
    history = [];
    historyIndex = -1;
    updateUndoRedoButtons();
}

// ============================================
// ZOOM
// ============================================
function applyZoom() {
    canvas.style.transform = `scale(${zoom})`;
    zoomLevelEl.textContent = Math.round(zoom * 100) + '%';
    updateConnections();
}

function zoomIn() {
    if (zoom < MAX_ZOOM) {
        zoom = Math.min(MAX_ZOOM, zoom + ZOOM_STEP);
        applyZoom();
    }
}

function zoomOut() {
    if (zoom > MIN_ZOOM) {
        zoom = Math.max(MIN_ZOOM, zoom - ZOOM_STEP);
        applyZoom();
    }
}

function resetZoom() {
    zoom = 1.0;
    applyZoom();
}

// Mouse wheel zoom (Ctrl + wheel)
canvasContainer.addEventListener('wheel', (e) => {
    if (e.ctrlKey || e.metaKey) {
        e.preventDefault();
        const delta = e.deltaY > 0 ? -ZOOM_STEP : ZOOM_STEP;
        const newZoom = Math.max(MIN_ZOOM, Math.min(MAX_ZOOM, zoom + delta));
        if (newZoom !== zoom) {
            zoom = newZoom;
            applyZoom();
        }
    }
}, {passive: false});

// Keyboard shortcuts
document.addEventListener('keydown', (e) => {
    if (e.ctrlKey || e.metaKey) {
        if (e.key === '=' || e.key === '+') {
            e.preventDefault();
            zoomIn();
        } else if (e.key === '-') {
            e.preventDefault();
            zoomOut();
        } else if (e.key === '0') {
            e.preventDefault();
            resetZoom();
        } else if (e.key === 'z' && !e.shiftKey) {
            e.preventDefault();
            undo();
        } else if ((e.key === 'y') || (e.key === 'z' && e.shiftKey)) {
            e.preventDefault();
            redo();
        }
    }
});

// ============================================
// COORDINATE CONVERSION (screen <-> canvas)
// ============================================
function screenToCanvas(sx, sy) {
    return {
        x: sx / zoom,
        y: sy / zoom
    };
}

function canvasToScreen(cx, cy) {
    return {
        x: cx * zoom,
        y: cy * zoom
    };
}

// ============================================
// BLOCK CREATION
// ============================================
function addBlock(type, data = null) {
    emptyState.style.display = 'none';

    const id = data ? data.id : nextId++;
    if (!data) nextId = Math.max(nextId, id + 1);

    const blockData = data || {
        id: id,
        type: type,
        name: type === 'interface' ? 'Reader' : type === 'struct' ? 'FileReader' : 'Read',
        x: 100 + (blocks.length * 30) % 400,
        y: 100 + (blocks.length * 20) % 300,
        width: 280,
        height: null,
        color: DEFAULT_COLORS[type],
        fields: [],
        methods: []
    };

    if (!data) {
        if (type === 'interface') {
            blockData.methods = [
                {name: 'Read', params: [{name: 'p', type: '[]byte'}], returns: ['int', 'error']}
            ];
        } else if (type === 'struct') {
            blockData.fields = [
                {name: 'data', type: '[]byte'},
                {name: 'pos', type: 'int'}
            ];
            blockData.methods = [
                {name: 'Read', params: [{name: 'p', type: '[]byte'}], returns: ['int', 'error']}
            ];
        } else if (type === 'method') {
            blockData.methods = [
                {name: 'Process', params: [{name: 'input', type: 'string'}], returns: ['string', 'error']}
            ];
        }
    }

    blocks.push(blockData);
    renderBlock(blockData);
    updateConnections();
    saveState();
}

function renderBlock(block) {
    const el = document.createElement('div');
    el.className = `block ${block.type}`;
    el.id = `block-${block.id}`;
    el.style.left = block.x + 'px';
    el.style.top = block.y + 'px';
    if (block.width) el.style.width = block.width + 'px';
    if (block.height) el.style.height = block.height + 'px';
    el.style.minWidth = '240px';
    el.style.minHeight = '80px';
    el.style.borderColor = block.color || DEFAULT_COLORS[block.type];

    let icon = block.type === 'interface' ? 'I' : block.type === 'struct' ? 'S' : 'F';
    let typeLabel = block.type === 'interface' ? 'interface' : block.type === 'struct' ? 'struct' : 'func';

    const colorPickerHtml = `
                <div class="color-picker-btn" style="background:${block.color || DEFAULT_COLORS[block.type]}" 
                     onclick="toggleColorPicker(event, ${block.id})" data-block-id="${block.id}">
                    <div class="color-picker-dropdown" id="color-picker-${block.id}">
                        ${COLOR_PALETTE.map(c => `<div class="color-option ${c === (block.color || DEFAULT_COLORS[block.type]) ? 'selected' : ''}" 
                            style="background:${c}" onclick="setBlockColor(${block.id}, '${c}')"></div>`).join('')}
                    </div>
                </div>
            `;

    let methodsHtml = '';
    if (block.methods && block.methods.length > 0) {
        methodsHtml = `
                    <div class="section-title">
                        <span>Methods</span>
                        <span class="add-btn" onclick="addMethod(${block.id})">+</span>
                    </div>
                `;
        block.methods.forEach((method, idx) => {
            methodsHtml += renderMethodView(block.id, idx, method);
        });
    } else if (block.type !== 'method') {
        methodsHtml = `
                    <div class="section-title">
                        <span>Methods</span>
                        <span class="add-btn" onclick="addMethod(${block.id})">+</span>
                    </div>
                `;
    }

    let fieldsHtml = '';
    if (block.fields && block.fields.length > 0) {
        fieldsHtml = `
                    <div class="section-title">
                        <span>Fields</span>
                        <span class="add-btn" onclick="addField(${block.id})">+</span>
                    </div>
                `;
        block.fields.forEach((field, idx) => {
            fieldsHtml += `
                        <div class="field-row">
                            <input class="field-name" value="${escapeHtml(field.name)}" 
                                onchange="updateFieldName(${block.id}, ${idx}, this.value)">
                            <input class="field-type" value="${escapeHtml(field.type)}" 
                                onchange="updateFieldType(${block.id}, ${idx}, this.value)">
                            <span class="remove-field" onclick="removeField(${block.id}, ${idx})">×</span>
                        </div>
                    `;
        });
    } else if (block.type === 'struct') {
        fieldsHtml = `
                    <div class="section-title">
                        <span>Fields</span>
                        <span class="add-btn" onclick="addField(${block.id})">+</span>
                    </div>
                `;
    }

    el.innerHTML = `
                <div class="block-header">
                    <div class="block-icon" style="background:rgba(${hexToRgb(block.color || DEFAULT_COLORS[block.type])},0.15);color:${block.color || DEFAULT_COLORS[block.type]}">${icon}</div>
                    <input class="block-name-input" value="${escapeHtml(block.name)}" 
                        onchange="updateBlockName(${block.id}, this.value)">
                    <span style="color:#8b949e;font-size:10px;">${typeLabel}</span>
                    <div class="block-actions">
                        ${colorPickerHtml}
                        <button onclick="deleteBlock(${block.id})">×</button>
                    </div>
                </div>
                <div class="block-body" id="block-body-${block.id}">
                    ${fieldsHtml}
                    ${methodsHtml}
                </div>
                <div class="resize-handle" onmousedown="startResize(event, ${block.id})"></div>
            `;

    el.addEventListener('mousedown', (e) => startDrag(e, block.id));
    canvas.appendChild(el);
}

function renderMethodView(blockId, idx, method) {
    const fullSig = buildFullSignature(method);
    const sigString = methodToString(method);

    return `
                <div class="method-row" id="method-row-${blockId}-${idx}">
                    <div class="method-signature-display" id="method-view-${blockId}-${idx}">
                        <div class="method-full-signature">${fullSig}</div>
                        <div class="edit-actions">
                            <button class="edit-btn" onclick="startEditMethod(${blockId}, ${idx})" title="Edit">
                                ✏️
                            </button>
                            <button class="edit-btn delete" onclick="removeMethod(${blockId}, ${idx})" title="Delete">
                                ×
                            </button>
                        </div>
                    </div>
                    <div class="method-signature-edit" id="method-edit-${blockId}-${idx}" style="display:none;width:100%">
                        <input class="method-edit-input" value="${escapeHtml(sigString)}" 
                            id="method-input-${blockId}-${idx}"
                            onkeydown="handleMethodEditKey(event, ${blockId}, ${idx})">
                        <div class="edit-actions">
                            <button class="edit-btn save" onclick="saveMethodEdit(${blockId}, ${idx})" title="Save (Enter)">
                                ✓
                            </button>
                            <button class="edit-btn cancel" onclick="cancelMethodEdit(${blockId}, ${idx})" title="Undo (Esc)">
                                ✕
                            </button>
                        </div>
                    </div>
                </div>
            `;
}

function buildFullSignature(method) {
    let html = `<span class="method-name">${escapeHtml(method.name)}</span>`;
    html += `<span class="method-paren">(</span>`;
    method.params.forEach((p, i) => {
        if (i > 0) html += '<span class="method-comma">, </span>';
        html += `<span class="method-param-name">${escapeHtml(p.name)}</span>`;
        html += `<span class="method-paren"> </span>`;
        html += `<span class="method-param-type">${escapeHtml(p.type)}</span>`;
    });
    html += `<span class="method-paren">)</span>`;
    if (method.returns.length > 0) {
        html += ' ';
        if (method.returns.length > 1) html += '<span class="method-paren">(</span>';
        method.returns.forEach((r, i) => {
            if (i > 0) html += '<span class="method-comma">, </span>';
            html += `<span class="method-return">${escapeHtml(r)}</span>`;
        });
        if (method.returns.length > 1) html += '<span class="method-paren">)</span>';
    }
    return html;
}

function methodToString(method) {
    const params = method.params.map(p => p.name + ' ' + p.type).join(', ');
    const returns = method.returns.join(', ');
    let s = `${method.name}(${params})`;
    if (returns) s += ` (${returns})`;
    return s;
}

function startEditMethod(blockId, idx) {
    const viewEl = document.getElementById(`method-view-${blockId}-${idx}`);
    const editEl = document.getElementById(`method-edit-${blockId}-${idx}`);
    const inputEl = document.getElementById(`method-input-${blockId}-${idx}`);

    if (viewEl && editEl && inputEl) {
        viewEl.style.display = 'none';
        editEl.style.display = 'flex';
        editEl.style.alignItems = 'center';
        editEl.style.gap = '6px';
        inputEl.focus();
        inputEl.select();
    }
}

function saveMethodEdit(blockId, idx) {
    const inputEl = document.getElementById(`method-input-${blockId}-${idx}`);
    if (!inputEl) return;

    const signature = inputEl.value;
    const block = blocks.find(b => b.id === blockId);
    const parsed = parseMethodSignature(signature);
    block.methods[idx] = parsed;

    const rowEl = document.getElementById(`method-row-${blockId}-${idx}`);
    if (rowEl) {
        rowEl.outerHTML = renderMethodView(blockId, idx, parsed);
    }

    updateConnections();
    saveState();
}

function cancelMethodEdit(blockId, idx) {
    const block = blocks.find(b => b.id === blockId);
    const method = block.methods[idx];

    const rowEl = document.getElementById(`method-row-${blockId}-${idx}`);
    if (rowEl) {
        rowEl.outerHTML = renderMethodView(blockId, idx, method);
    }
}

function handleMethodEditKey(event, blockId, idx) {
    if (event.key === 'Enter') {
        event.preventDefault();
        saveMethodEdit(blockId, idx);
    } else if (event.key === 'Escape') {
        event.preventDefault();
        cancelMethodEdit(blockId, idx);
    }
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function hexToRgb(hex) {
    const r = parseInt(hex.slice(1, 3), 16);
    const g = parseInt(hex.slice(3, 5), 16);
    const b = parseInt(hex.slice(5, 7), 16);
    return `${r},${g},${b}`;
}

// ============================================
// COLOR PICKER
// ============================================
function toggleColorPicker(event, blockId) {
    event.stopPropagation();
    const dropdown = document.getElementById(`color-picker-${blockId}`);
    const isActive = dropdown.classList.contains('active');
    document.querySelectorAll('.color-picker-dropdown').forEach(d => d.classList.remove('active'));
    if (!isActive) dropdown.classList.add('active');
}

function setBlockColor(blockId, color) {
    const block = blocks.find(b => b.id === blockId);
    if (block) {
        block.color = color;
        refreshBlock(block);
        updateConnections();
        saveState();
    }
    document.querySelectorAll('.color-picker-dropdown').forEach(d => d.classList.remove('active'));
}

document.addEventListener('click', (e) => {
    if (!e.target.closest('.color-picker-btn')) {
        document.querySelectorAll('.color-picker-dropdown').forEach(d => d.classList.remove('active'));
    }
});

// ============================================
// RESIZE
// ============================================
function startResize(event, blockId) {
    event.stopPropagation();
    event.preventDefault();
    resizingBlock = blocks.find(b => b.id === blockId);
    const el = document.getElementById(`block-${blockId}`);
    resizeStart = {
        x: event.clientX,
        y: event.clientY,
        w: el.offsetWidth,
        h: el.offsetHeight
    };
}

document.addEventListener('mousemove', (e) => {
    if (resizingBlock) {
        const dx = (e.clientX - resizeStart.x) / zoom;
        const dy = (e.clientY - resizeStart.y) / zoom;
        const newW = Math.max(240, resizeStart.w + dx);
        const newH = Math.max(80, resizeStart.h + dy);

        resizingBlock.width = newW;
        resizingBlock.height = newH;

        const el = document.getElementById(`block-${resizingBlock.id}`);
        el.style.width = newW + 'px';
        el.style.height = newH + 'px';

        updateConnections();
    }

    if (draggedBlock) {
        const containerRect = canvasContainer.getBoundingClientRect();
        const rawX = e.clientX - containerRect.left + canvasContainer.scrollLeft;
        const rawY = e.clientY - containerRect.top + canvasContainer.scrollTop;

        const newX = (rawX - dragOffset.x) / zoom;
        const newY = (rawY - dragOffset.y) / zoom;

        draggedBlock.x = Math.max(0, newX);
        draggedBlock.y = Math.max(0, newY);

        const el = document.getElementById(`block-${draggedBlock.id}`);
        el.style.left = draggedBlock.x + 'px';
        el.style.top = draggedBlock.y + 'px';

        updateConnections();
    }

    if (isPanning) {
        const dx = e.clientX - panStart.x;
        const dy = e.clientY - panStart.y;
        canvasContainer.scrollLeft -= dx / zoom;
        canvasContainer.scrollTop -= dy / zoom;
        panStart = {x: e.clientX, y: e.clientY};
    }
});

document.addEventListener('mouseup', () => {
    if (resizingBlock || draggedBlock) {
        saveState();
    }
    resizingBlock = null;
    draggedBlock = null;
    isPanning = false;
});

// ============================================
// DRAG & DROP
// ============================================
function startDrag(e, blockId) {
    if (e.target.tagName === 'INPUT' || e.target.tagName === 'BUTTON' ||
        e.target.classList.contains('add-btn') || e.target.classList.contains('remove-field') ||
        e.target.classList.contains('resize-handle') || e.target.closest('.color-picker-dropdown') ||
        e.target.classList.contains('color-picker-btn') || e.target.classList.contains('color-option') ||
        e.target.classList.contains('edit-btn') || e.target.closest('.edit-actions')) {
        return;
    }
    e.preventDefault();
    draggedBlock = blocks.find(b => b.id === blockId);
    const el = document.getElementById(`block-${blockId}`);
    const rect = el.getBoundingClientRect();
    dragOffset.x = (e.clientX - rect.left);
    dragOffset.y = (e.clientY - rect.top);
}

canvasContainer.addEventListener('mousedown', (e) => {
    if (e.target === canvasContainer || e.target === canvas) {
        isPanning = true;
        panStart = {x: e.clientX, y: e.clientY};
    }
});

// ============================================
// CONNECTIONS
// ============================================
function updateConnections() {
    document.querySelectorAll('.connection-line').forEach(el => el.remove());
    connections = [];
    document.querySelectorAll('.implementation-badge').forEach(el => el.remove());

    const interfaces = blocks.filter(b => b.type === 'interface');
    const structs = blocks.filter(b => b.type === 'struct');

    structs.forEach(struct => {
        interfaces.forEach(iface => {
            if (structImplementsInterface(struct, iface)) {
                drawConnection(struct, iface);
            }
        });
    });
}

function structImplementsInterface(struct, iface) {
    if (!iface.methods || iface.methods.length === 0) return false;
    if (!struct.methods || struct.methods.length === 0) return false;

    return iface.methods.every(imethod => {
        return struct.methods.some(smethod => {
            return methodsMatch(imethod, smethod);
        });
    });
}

function methodsMatch(m1, m2) {
    if (m1.name !== m2.name) return false;
    if (m1.params.length !== m2.params.length) return false;
    if (m1.returns.length !== m2.returns.length) return false;

    for (let i = 0; i < m1.params.length; i++) {
        if (m1.params[i].type !== m2.params[i].type) return false;
    }
    for (let i = 0; i < m1.returns.length; i++) {
        if (m1.returns[i] !== m2.returns[i]) return false;
    }
    return true;
}

function drawConnection(struct, iface) {
    const structEl = document.getElementById(`block-${struct.id}`);
    const ifaceEl = document.getElementById(`block-${iface.id}`);
    if (!structEl || !ifaceEl) return;

    const sRect = {left: struct.x, top: struct.y, width: structEl.offsetWidth, height: structEl.offsetHeight};
    const iRect = {left: iface.x, top: iface.y, width: ifaceEl.offsetWidth, height: ifaceEl.offsetHeight};

    const startX = sRect.left + sRect.width / 2;
    const startY = sRect.top;
    const endX = iRect.left + iRect.width / 2;
    const endY = iRect.top + iRect.height;

    const line = document.createElement('div');
    line.className = 'connection-line';

    const svgNS = 'http://www.w3.org/2000/svg';
    const svg = document.createElementNS(svgNS, 'svg');
    svg.style.position = 'absolute';
    svg.style.left = '0';
    svg.style.top = '0';
    svg.style.width = '5000px';
    svg.style.height = '5000px';
    svg.style.pointerEvents = 'none';

    const path = document.createElementNS(svgNS, 'path');
    const midY = (startY + endY) / 2;
    const d = `M ${startX} ${startY} C ${startX} ${midY}, ${endX} ${midY}, ${endX} ${endY}`;
    path.setAttribute('d', d);
    path.setAttribute('class', 'connection-path');
    path.style.stroke = iface.color || DEFAULT_COLORS.interface;

    svg.appendChild(path);
    line.appendChild(svg);
    canvas.appendChild(line);

    connections.push({struct: struct.id, interface: iface.id});

    const badge = document.createElement('div');
    badge.className = 'implementation-badge';
    badge.textContent = 'implements';
    badge.title = `Implements ${iface.name}`;
    badge.style.background = iface.color || DEFAULT_COLORS.interface;
    structEl.appendChild(badge);
}

// ============================================
// BLOCK EDITING
// ============================================
function updateBlockName(id, name) {
    const block = blocks.find(b => b.id === id);
    if (block) block.name = name;
    updateConnections();
    saveState();
}

function deleteBlock(id) {
    blocks = blocks.filter(b => b.id !== id);
    const el = document.getElementById(`block-${id}`);
    if (el) el.remove();
    if (blocks.length === 0) emptyState.style.display = 'block';
    updateConnections();
    saveState();
}

function addField(blockId) {
    const block = blocks.find(b => b.id === blockId);
    if (!block.fields) block.fields = [];
    block.fields.push({name: 'newField', type: 'string'});
    refreshBlock(block);
    updateConnections();
    saveState();
}

function removeField(blockId, idx) {
    const block = blocks.find(b => b.id === blockId);
    block.fields.splice(idx, 1);
    refreshBlock(block);
    updateConnections();
    saveState();
}

function updateFieldName(blockId, idx, name) {
    const block = blocks.find(b => b.id === blockId);
    block.fields[idx].name = name;
    saveState();
}

function updateFieldType(blockId, idx, type) {
    const block = blocks.find(b => b.id === blockId);
    block.fields[idx].type = type;
    saveState();
}

function addMethod(blockId) {
    const block = blocks.find(b => b.id === blockId);
    if (!block.methods) block.methods = [];
    block.methods.push({name: 'NewMethod', params: [], returns: []});
    refreshBlock(block);
    updateConnections();
    saveState();
}

function removeMethod(blockId, idx) {
    const block = blocks.find(b => b.id === blockId);
    block.methods.splice(idx, 1);
    refreshBlock(block);
    updateConnections();
    saveState();
}

function parseMethodSignature(sig) {
    sig = sig.trim();
    const firstParen = sig.indexOf('(');
    if (firstParen === -1) {
        return {name: sig || 'Method', params: [], returns: []};
    }

    const name = sig.substring(0, firstParen).trim();

    let parenDepth = 0;
    let paramsEnd = -1;
    for (let i = firstParen; i < sig.length; i++) {
        if (sig[i] === '(') parenDepth++;
        if (sig[i] === ')') {
            parenDepth--;
            if (parenDepth === 0) {
                paramsEnd = i;
                break;
            }
        }
    }

    if (paramsEnd === -1) {
        return {name: name || 'Method', params: [], returns: []};
    }

    const paramsStr = sig.substring(firstParen + 1, paramsEnd);
    const params = parseParams(paramsStr);

    let returnsStr = sig.substring(paramsEnd + 1).trim();
    const returns = parseReturns(returnsStr);

    return {name: name || 'Method', params, returns};
}

function parseParams(paramsStr) {
    const params = [];
    if (!paramsStr.trim()) return params;

    const parts = splitByComma(paramsStr);

    parts.forEach(part => {
        const trimmed = part.trim();
        if (!trimmed) return;

        const spaceIdx = findNameTypeSplit(trimmed);

        if (spaceIdx > 0) {
            params.push({
                name: trimmed.substring(0, spaceIdx).trim(),
                type: trimmed.substring(spaceIdx + 1).trim()
            });
        } else {
            params.push({name: trimmed, type: 'interface{}'});
        }
    });

    return params;
}

function findNameTypeSplit(str) {
    let depth = 0;
    for (let i = str.length - 1; i >= 0; i--) {
        const ch = str[i];
        if (ch === ']' || ch === ')') depth++;
        if (ch === '[' || ch === '(') depth--;
        if (ch === ' ' && depth === 0) {
            return i;
        }
    }
    return -1;
}

function splitByComma(str) {
    const parts = [];
    let current = '';
    let depth = 0;

    for (let i = 0; i < str.length; i++) {
        const ch = str[i];
        if (ch === '(' || ch === '[' || ch === '{') depth++;
        if (ch === ')' || ch === ']' || ch === '}') depth--;

        if (ch === ',' && depth === 0) {
            parts.push(current);
            current = '';
        } else {
            current += ch;
        }
    }
    if (current.trim()) parts.push(current);
    return parts;
}

function parseReturns(returnsStr) {
    const returns = [];
    if (!returnsStr) return returns;

    returnsStr = returnsStr.trim();
    if (returnsStr.startsWith('(') && returnsStr.endsWith(')')) {
        returnsStr = returnsStr.substring(1, returnsStr.length - 1);
    }

    if (!returnsStr.trim()) return returns;

    const parts = splitByComma(returnsStr);
    parts.forEach(part => {
        const trimmed = part.trim();
        if (trimmed) returns.push(trimmed);
    });

    return returns;
}

function refreshBlock(block) {
    const el = document.getElementById(`block-${block.id}`);
    if (el) {
        // Reset height to auto so content fits, but keep width
        block.height = null;
        el.remove();
        renderBlock(block);
        updateConnections();
    }
}

function clearCanvas() {
    blocks = [];
    connections = [];
    document.querySelectorAll('.block').forEach(el => el.remove());
    document.querySelectorAll('.connection-line').forEach(el => el.remove());
    emptyState.style.display = 'block';
    saveState();
}

// ============================================
// EXPORT / IMPORT
// ============================================
function exportProject() {
    let code = 'package main\n\n';

    blocks.filter(b => b.type === 'interface').forEach(block => {
        code += `// ${block.name} interface\n`;
        code += `type ${block.name} interface {\n`;
        block.methods.forEach(m => {
            const params = m.params.map(p => `${p.name} ${p.type}`).join(', ');
            const returns = m.returns.join(', ');
            code += `    ${m.name}(${params})`;
            if (returns) code += ` (${returns})`;
            code += `\n`;
        });
        code += `}\n\n`;
    });

    blocks.filter(b => b.type === 'struct').forEach(block => {
        code += `// ${block.name} struct\n`;
        code += `type ${block.name} struct {\n`;
        block.fields.forEach(f => {
            code += `    ${f.name} ${f.type}\n`;
        });
        code += `}\n\n`;

        block.methods.forEach(m => {
            const params = m.params.map(p => `${p.name} ${p.type}`).join(', ');
            const returns = m.returns.join(', ');
            code += `func (s *${block.name}) ${m.name}(${params})`;
            if (returns) code += ` (${returns})`;
            code += ` {\n`;
            code += `    // TODO: implement\n`;
            if (m.returns.length > 0) {
                const zeroReturns = m.returns.map(r => {
                    if (r === 'error') return 'nil';
                    if (r === 'string') return '""';
                    if (r === 'int' || r === 'int64' || r === 'float64') return '0';
                    if (r === 'bool') return 'false';
                    return 'nil';
                }).join(', ');
                code += `    return ${zeroReturns}\n`;
            }
            code += `}\n\n`;
        });
    });

    blocks.filter(b => b.type === 'method').forEach(block => {
        block.methods.forEach(m => {
            const params = m.params.map(p => `${p.name} ${p.type}`).join(', ');
            const returns = m.returns.join(', ');
            code += `func ${m.name}(${params})`;
            if (returns) code += ` (${returns})`;
            code += ` {\n`;
            code += `    // TODO: implement\n`;
            if (m.returns.length > 0) {
                const zeroReturns = m.returns.map(r => {
                    if (r === 'error') return 'nil';
                    if (r === 'string') return '""';
                    if (r === 'int' || r === 'int64' || r === 'float64') return '0';
                    if (r === 'bool') return 'false';
                    return 'nil';
                }).join(', ');
                code += `    return ${zeroReturns}\n`;
            }
            code += `}\n\n`;
        });
    });

    connections.forEach(conn => {
        const struct = blocks.find(b => b.id === conn.struct);
        const iface = blocks.find(b => b.id === conn.interface);
        if (struct && iface) {
            code += `// Compile-time interface check\n`;
            code += `var _ ${iface.name} = (*${struct.name})(nil)\n\n`;
        }
    });

    showModal('Export Go Code', code, 'Copy', () => {
        navigator.clipboard.writeText(code);
        document.getElementById('modalActionBtn').textContent = 'Copied!';
        setTimeout(() => {
            document.getElementById('modalActionBtn').textContent = 'Copy';
        }, 2000);
    });
}

function importProject() {
    showModal('Import Go Code', '', 'Import', () => {
        const code = document.getElementById('modalTextarea').value;
        parseGoCode(code);
        closeModal();
    });
}

function parseGoCode(code) {
    clearCanvas();

    const interfaceRegex = /type\s+(\w+)\s+interface\s*\{([^}]*)\}/g;
    let match;
    while ((match = interfaceRegex.exec(code)) !== null) {
        const name = match[1];
        const body = match[2];
        const methods = [];

        const methodLines = body.split('\n').filter(l => l.trim());
        methodLines.forEach(line => {
            const methodMatch = line.trim().match(/^(\w+)\s*\((.*?)\)\s*(?:\((.*?)\)|(.*))?$/);
            if (methodMatch) {
                const mName = methodMatch[1];
                const paramsStr = methodMatch[2];
                const returnsStr = methodMatch[3] || methodMatch[4] || '';

                const params = parseParams(paramsStr);
                const returns = parseReturns(returnsStr);

                methods.push({name: mName, params, returns});
            }
        });

        addBlock('interface', {
            id: nextId++,
            type: 'interface',
            name: name,
            x: 100 + (blocks.length * 50),
            y: 100,
            fields: [],
            methods: methods
        });
    }

    const structRegex = /type\s+(\w+)\s+struct\s*\{([^}]*)\}/g;
    while ((match = structRegex.exec(code)) !== null) {
        const name = match[1];
        const body = match[2];
        const fields = [];

        const fieldLines = body.split('\n').filter(l => l.trim());
        fieldLines.forEach(line => {
            const fieldMatch = line.trim().match(/^(\w+)\s+(\S+(?:\[.*?\])?(?:\{.*?\})?)$/);
            if (fieldMatch) {
                fields.push({name: fieldMatch[1], type: fieldMatch[2]});
            }
        });

        const methods = [];
        const methodRegex = new RegExp(`func\\s*\\(s\\s*\\*?${name}\\)\\s+(\\w+)\\s*\\((.*?)\\)\\s*(?:\\((.*?)\\)|(.*))?`, 'g');
        let mMatch;
        while ((mMatch = methodRegex.exec(code)) !== null) {
            const mName = mMatch[1];
            const paramsStr = mMatch[2];
            const returnsStr = mMatch[3] || mMatch[4] || '';

            const params = parseParams(paramsStr);
            const returns = parseReturns(returnsStr);

            methods.push({name: mName, params, returns});
        }

        addBlock('struct', {
            id: nextId++,
            type: 'struct',
            name: name,
            x: 100 + (blocks.length * 50),
            y: 300,
            fields: fields,
            methods: methods
        });
    }

    const funcRegex = /func\s+(\w+)\s*\((.*?)\)\s*(?:\((.*?)\)|(.*))?\s*\{/g;
    while ((match = funcRegex.exec(code)) !== null) {
        const name = match[1];
        const paramsStr = match[2];
        const returnsStr = match[3] || match[4] || '';

        const params = parseParams(paramsStr);
        const returns = parseReturns(returnsStr);

        const prevCode = code.substring(Math.max(0, match.index - 20), match.index);
        if (!prevCode.includes(')')) {
            addBlock('method', {
                id: nextId++,
                type: 'method',
                name: name,
                x: 100 + (blocks.length * 50),
                y: 500,
                fields: [],
                methods: [{name: name, params, returns}]
            });
        }
    }

    updateConnections();
}

// ============================================
// MODAL
// ============================================
function showModal(title, text, actionText, action) {
    document.getElementById('modalTitle').textContent = title;
    document.getElementById('modalTextarea').value = text;
    const btn = document.getElementById('modalActionBtn');
    btn.textContent = actionText;
    btn.onclick = action;
    document.getElementById('modalOverlay').classList.add('active');
}

function closeModal() {
    document.getElementById('modalOverlay').classList.remove('active');
}

function modalAction() {
    // Set dynamically
}

document.getElementById('modalOverlay').addEventListener('click', (e) => {
    if (e.target === document.getElementById('modalOverlay')) closeModal();
});

// ============================================
// DEMO DATA
// ============================================
window.addEventListener('load', () => {
    // Try to load from localStorage first
    if (loadFromLocalStorage()) {
        // Restore canvas from loaded state
        blocks.forEach(block => renderBlock(block));
        if (blocks.length > 0) {
            emptyState.style.display = 'none';
        }
        applyZoom();
        updateConnections();
        updateUndoRedoButtons();

        // Initialize history if empty
        if (history.length === 0) {
            saveState();
        }
        return;
    }

    // Initialize empty history
    clearHistory();

    addBlock('interface', {
        id: nextId++,
        type: 'interface',
        name: 'Reader',
        x: 100, y: 80,
        width: 320,
        color: '#58a6ff',
        fields: [],
        methods: [
            {name: 'Read', params: [{name: 'p', type: '[]byte'}], returns: ['int', 'error']}
        ]
    });

    addBlock('interface', {
        id: nextId++,
        type: 'interface',
        name: 'Writer',
        x: 500, y: 80,
        width: 320,
        color: '#a371f7',
        fields: [],
        methods: [
            {name: 'Write', params: [{name: 'p', type: '[]byte'}], returns: ['int', 'error']}
        ]
    });

    addBlock('struct', {
        id: nextId++,
        type: 'struct',
        name: 'FileReader',
        x: 100, y: 350,
        width: 340,
        color: '#3fb950',
        fields: [
            {name: 'file', type: '*os.File'},
            {name: 'buffer', type: '[]byte'}
        ],
        methods: [
            {name: 'Read', params: [{name: 'p', type: '[]byte'}], returns: ['int', 'error']}
        ]
    });

    addBlock('struct', {
        id: nextId++,
        type: 'struct',
        name: 'Buffer',
        x: 500, y: 350,
        width: 340,
        color: '#d29922',
        fields: [
            {name: 'data', type: '[]byte'},
            {name: 'off', type: 'int'}
        ],
        methods: [
            {name: 'Read', params: [{name: 'p', type: '[]byte'}], returns: ['int', 'error']},
            {name: 'Write', params: [{name: 'p', type: '[]byte'}], returns: ['int', 'error']}
        ]
    });

    addBlock('method', {
        id: nextId++,
        type: 'method',
        name: 'ProcessData',
        x: 900, y: 200,
        width: 400,
        color: '#f85149',
        fields: [],
        methods: [
            {
                name: 'ProcessData',
                params: [{name: 'input', type: 'string'}, {name: 'opts', type: 'Options'}],
                returns: ['string', 'error']
            }
        ]
    });
    updateConnections();
});