document.addEventListener('DOMContentLoaded', function() {
    // Инициализация карты
    const ctx = document.getElementById('map-chart').getContext('2d');
    const chart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: [],
            datasets: [{
                label: 'Данные по регионам',
                data: [],
                backgroundColor: 'rgba(54, 162, 235, 0.5)',
                borderColor: 'rgba(54, 162, 235, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        }
    });

    // Обработчики вкладок
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.addEventListener('click', function() {
            document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
            document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));

            this.classList.add('active');
            document.getElementById(this.dataset.tab).classList.add('active');
        });
    });

    // Обработчик кнопки "Применить"
    document.getElementById('apply-filters').addEventListener('click', function() {
        const region = document.getElementById('region').value;
        const dateFrom = document.getElementById('date_from').value;
        const dateTo = document.getElementById('date_to').value;
        const searchText = document.getElementById('search-input').value;
        const extension = document.getElementById('extension').value;

        // Загрузка данных для карты и статистики
        loadDataForChart(region, dateFrom, dateTo, chart);

        // Загрузка файлов
        if (searchText) {
            searchFiles(region, searchText, extension);
        } else {
            loadFilesForRegion(region);
        }
    });

    // Первоначальная загрузка данных
    loadDataForChart('all', '', '', chart);
});

function loadDataForChart(region, dateFrom, dateTo, chart) {
    // Здесь остается ваша существующая логика загрузки данных для карты
    // ...
}

function searchFiles(region, searchText, extension) {
    fetch(`/api/search?region=${region}&q=${encodeURIComponent(searchText)}&ext=${extension}`)
        .then(response => response.json())
        .then(files => {
            updateFilesTable(files);
        })
        .catch(error => console.error('Ошибка поиска:', error));
}

function loadFilesForRegion(region) {
    fetch(`/api/files?region=${region}`)
        .then(response => response.json())
        .then(files => {
            updateFilesTable(files);
        })
        .catch(error => console.error('Ошибка загрузки файлов:', error));
}

function updateFilesTable(files) {
    const tbody = document.getElementById('files-body');
    tbody.innerHTML = '';

    files.forEach(file => {
        const row = document.createElement('tr');

        const nameCell = document.createElement('td');
        nameCell.textContent = file.Filename || file.name;
        row.appendChild(nameCell);

        const pathCell = document.createElement('td');
        pathCell.textContent = file.Path || file.path;
        row.appendChild(pathCell);

        const regionCell = document.createElement('td');
        regionCell.textContent = file.Region || 'Не указан';
        row.appendChild(regionCell);

        const actionCell = document.createElement('td');
        const link = document.createElement('a');
        link.href = file.Path || file.path;
        link.textContent = 'Открыть';
        link.target = '_blank';
        actionCell.appendChild(link);
        row.appendChild(actionCell);

        tbody.appendChild(row);
    });
}