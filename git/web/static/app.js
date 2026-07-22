const API_BASE = '/api';

const IncomeCategories = ['工资', '奖金', '投资收益', '兼职', '其他收入'];
const ExpenseCategories = ['餐饮', '交通', '购物', '娱乐', '医疗', '教育', '住房', '其他支出'];

let allRecords = [];
let currentUser = null;

function checkLoginStatus() {
    const savedUser = sessionStorage.getItem('currentUser');
    if (savedUser) {
        currentUser = JSON.parse(savedUser);
        document.getElementById('auth-section').classList.add('hidden');
        document.getElementById('dashboard-section').classList.remove('hidden');
        document.getElementById('current-user').textContent = `欢迎, ${currentUser.name}`;
        switchView('add');
    }
}

window.addEventListener('load', checkLoginStatus);

function switchTab(tab) {
    document.getElementById('register-form').classList.add('hidden');
    document.getElementById('login-form').classList.add('hidden');
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    
    if (tab === 'register') {
        document.getElementById('register-form').classList.remove('hidden');
        document.querySelectorAll('.tab-btn')[0].classList.add('active');
    } else {
        document.getElementById('login-form').classList.remove('hidden');
        document.querySelectorAll('.tab-btn')[1].classList.add('active');
    }
    document.getElementById('auth-message').textContent = '';
}

function switchView(view) {
    document.getElementById('add-view').classList.add('hidden');
    document.getElementById('list-view').classList.add('hidden');
    document.querySelectorAll('.nav-btn').forEach(btn => btn.classList.remove('active'));
    
    if (view === 'add') {
        document.getElementById('add-view').classList.remove('hidden');
        document.querySelectorAll('.nav-btn')[0].classList.add('active');
    } else {
        document.getElementById('list-view').classList.remove('hidden');
        document.querySelectorAll('.nav-btn')[1].classList.add('active');
        loadRecords();
    }
}

function updateCategorySelect() {
    const sort = document.querySelector('input[name="sort"]:checked').value;
    const select = document.getElementById('category-select');
    const categories = sort === 'Income' ? IncomeCategories : ExpenseCategories;
    
    select.innerHTML = '<option value="">请选择分类</option>';
    categories.forEach(cat => {
        const option = document.createElement('option');
        option.value = cat;
        option.textContent = cat;
        select.appendChild(option);
    });
}

document.querySelectorAll('input[name="sort"]').forEach(radio => {
    radio.addEventListener('change', updateCategorySelect);
});

updateCategorySelect();

async function register() {
    const name = document.getElementById('reg-name').value;
    const password = document.getElementById('reg-password').value;
    
    if (!name || !password) {
        showMessage('auth-message', '请填写用户名和密码');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: name, password })
        });
        
        const data = await response.json();
        
        if (data.success) {
            showMessage('auth-message', data.message, true);
            setTimeout(() => {
                switchTab('login');
            }, 1500);
        } else {
            showMessage('auth-message', data.message);
        }
    } catch (error) {
        showMessage('auth-message', '注册失败，请检查网络');
    }
}

async function login() {
    const name = document.getElementById('login-name').value;
    const password = document.getElementById('login-password').value;
    
    if (!name || !password) {
        showMessage('auth-message', '请填写用户名和密码');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ username: name, password })
        });
        
        const data = await response.json();
        
        if (data.success && data.data && data.data.name) {
            currentUser = data.data;
            sessionStorage.setItem('currentUser', JSON.stringify(currentUser));
            document.getElementById('auth-section').classList.add('hidden');
            document.getElementById('dashboard-section').classList.remove('hidden');
            document.getElementById('current-user').textContent = `欢迎, ${data.data.name}`;
            switchView('add');
        } else if (data.success) {
            showMessage('auth-message', '登录成功，但用户信息不完整');
        } else {
            showMessage('auth-message', data.message || '登录失败');
        }
    } catch (error) {
        showMessage('auth-message', '登录失败，请检查网络');
    }
}

function logout() {
    currentUser = null;
    sessionStorage.removeItem('currentUser');
    document.getElementById('dashboard-section').classList.add('hidden');
    document.getElementById('auth-section').classList.remove('hidden');
    switchTab('login');
    document.getElementById('login-name').value = '';
    document.getElementById('login-password').value = '';
}

async function createRecord() {
    const sort = document.querySelector('input[name="sort"]:checked').value;
    const category = document.getElementById('category-select').value;
    const amount = parseFloat(document.getElementById('amount-input').value);
    const note = document.getElementById('note-input').value;
    
    if (!category) {
        showMessage('add-message', '请选择分类');
        return;
    }
    if (!amount || amount <= 0) {
        showMessage('add-message', '请输入有效金额');
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/CreateRecord`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_id: currentUser.id, sort, category, amount, note, date: new Date().toISOString(), total: amount })
        });
        
        const data = await response.json();
        
        if (data.success) {
            showMessage('add-message', data.message, true);
            document.getElementById('category-select').value = '';
            document.getElementById('amount-input').value = '';
            document.getElementById('note-input').value = '';
            allRecords.push(data.data);
        } else {
            showMessage('add-message', data.message);
        }
    } catch (error) {
        showMessage('add-message', '创建失败，请检查网�?);
    }
}

async function loadRecords() {
    try {
        const response = await fetch(`${API_BASE}/ShowRecord1`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ user_id: currentUser.id })
        });
        const data = await response.json();
        
        if (data.success) {
            allRecords = data.data;
            renderRecords(allRecords);
        } else {
            showMessage('list-message', data.message);
        }
    } catch (error) {
        showMessage('list-message', '加载失败，请检查网络');
    }
}

function filterRecords() {
    const searchText = document.getElementById('search-input').value.toLowerCase();
    const sortFilter = document.getElementById('filter-sort').value;
    
    let filtered = allRecords;
    
    if (sortFilter) {
        filtered = filtered.filter(r => r.sort === sortFilter);
    }
    
    if (searchText) {
        filtered = filtered.filter(r => r.note.toLowerCase().includes(searchText));
    }
    
    renderRecords(filtered);
}

function renderRecords(records) {
    const tbody = document.getElementById('records-table').querySelector('tbody');
    tbody.innerHTML = '';
    
    if (records.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" style="text-align: center;">暂无记录</td></tr>';
        return;
    }
    
    records.forEach(record => {
        const tr = document.createElement('tr');
        const sortText = record.sort === 'Income' ? '收入' : '支出';
        const sortClass = record.sort === 'Income' ? 'sort-income' : 'sort-expense';
        
        tr.innerHTML = `
            <td>${record.id}</td>
            <td><span class="${sortClass}">${sortText}</span></td>
            <td>${record.category}</td>
            <td>${record.amount.toFixed(2)}</td>
            <td>${record.note || '-'}</td>
            <td>${formatDate(record.date)}</td>
            <td><button class="delete-btn" onclick="deleteRecord(${record.id})">删除</button></td>
        `;
        
        tbody.appendChild(tr);
    });
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
}

async function deleteRecord(id) {
    if (!confirm('确定要删除这条记录吗？')) {
        return;
    }
    
    try {
        const response = await fetch(`${API_BASE}/DeleteRecord1`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id, user_id: currentUser.id })
        });
        
        const data = await response.json();
        
        if (data.success) {
            loadRecords();
            showMessage('list-message', '删除成功', true);
            setTimeout(() => showMessage('list-message', ''), 2000);
        } else {
            showMessage('list-message', data.message);
        }
    } catch (error) {
        showMessage('list-message', '删除失败，请检查网络');
    }
}

function showMessage(elementId, message, isSuccess = false) {
    const element = document.getElementById(elementId);
    element.textContent = message;
    element.className = 'message' + (isSuccess ? ' success' : '');
}
