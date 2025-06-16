// home_script.js
document.addEventListener('DOMContentLoaded', function() {
    // Инициализация карты (заглушка)
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

    // Обработчик кнопки "Применить"
    document.getElementById('apply-filters').addEventListener('click', function() {
        applyFilters(chart);
    });

    // Первоначальная загрузка данных
    applyFilters(chart);
});

function applyFilters(chart) {
    const region = document.getElementById('region').value;
    const dateFrom = document.getElementById('date_from').value;
    const dateTo = document.getElementById('date_to').value;

    // Формируем URL с параметрами
    let url = '/api/data?';
    if (region) url += `region=${region}&`;
    if (dateFrom) url += `date_from=${dateFrom}&`;
    if (dateTo) url += `date_to=${dateTo}&`;
    url = url.slice(0, -1); // Удаляем последний &

    // Запрашиваем данные с сервера
    fetch(url)
        .then(response => response.json())
        .then(data => {
            updateChart(chart, data);
            updateTable(data);
        })
        .catch(error => console.error('Ошибка:', error));
}

function updateChart(chart, data) {
    // Группируем данные по регионам
    const regions = {};
    data.forEach(item => {
        if (!regions[item.region]) {
            regions[item.region] = 0;
        }
        regions[item.region] += item.data;
    });

    // Обновляем данные карты
    chart.data.labels = Object.keys(regions);
    chart.data.datasets[0].data = Object.values(regions);
    chart.update();
}

function updateTable(data) {
    const tbody = document.getElementById('data-body');
    tbody.innerHTML = '';

    // Сортируем данные по дате (новые сначала)
    data.sort((a, b) => new Date(b.date) - new Date(a.date));

    // Заполняем таблицу
    data.forEach(item => {
        const row = document.createElement('tr');

        const regionCell = document.createElement('td');
        regionCell.textContent = getRegionName(item.region);
        row.appendChild(regionCell);

        const dataCell = document.createElement('td');
        dataCell.textContent = item.data;
        row.appendChild(dataCell);

        const dateCell = document.createElement('td');
        dateCell.textContent = new Date(item.date).toLocaleDateString();
        row.appendChild(dateCell);

        tbody.appendChild(row);
    });
}

function getRegionName(regionCode) {
    const regions = {
        'abay': 'Абайская область',
        'akmola': 'Акмолинская область',
        'aktobe': 'Актюбинская область',
        'almaty': 'Алматинская область',
        'atyrau': 'Атырауская область',
        'east-kazakhstan': 'Восточно-Казахстанская область',
        'zhambyl': 'Жамбылская область',
        'zhetysu': 'Жетысуская область',
        'west-kazakhstan': 'Западно-Казахстанская область',
        'karaganda': 'Карагандинская область',
        'kostanay': 'Костанайская область',
        'kzylorda': 'Кызылординская область',
        'mangystau': 'Мангистауская область',
        'pavlodar': 'Павлодарская область',
        'north-kazakhstan': 'Северо-Казахстанская область',
        'turkistan': 'Туркестанская область',
        'ulytau': 'Улытауская область',
        'almaty-city': 'г. Алматы',
        'astana': 'г. Астана',
        'shymkent': 'г. Шымкент'
    };
    return regions[regionCode] || regionCode;
}