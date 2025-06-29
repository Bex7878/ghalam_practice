document.addEventListener('DOMContentLoaded', function() {
    // Инициализация карты
    const map = L.map('map').setView([48.0, 67.0], 5); // Центр Казахстана

    // Добавляем слой OpenStreetMap
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);

    let markers = [];
    let regionLayers = {};

    // Обработчик кнопки "Применить"
    document.getElementById('apply-filters').addEventListener('click', function() {
        const region = document.getElementById('region').value;
        const dateFrom = document.getElementById('date_from').value;
        const dateTo = document.getElementById('date_to').value;
        const searchText = document.getElementById('search-input').value;
        const extension = document.getElementById('extension').value;

        searchFiles(region, dateFrom, dateTo, searchText, extension, map);
    });

    // Первоначальная загрузка данных
    searchFiles('all', '', '', '', '', map);
});

function clearMap(map) {
    // Удаляем все маркеры
    for (const region in regionLayers) {
        map.removeLayer(regionLayers[region]);
    }
    regionLayers = {};
}

async function searchFiles(region, dateFrom, dateTo, searchText, extension, map) {
    try {
        const params = new URLSearchParams();
        params.append('region', region);
        if (dateFrom) params.append('date_from', dateFrom);
        if (dateTo) params.append('date_to', dateTo);
        if (searchText) params.append('q', searchText);
        if (extension) params.append('ext', extension);

        const response = await fetch(`/api/search?${params.toString()}`);

        if (!response.ok) throw new Error('Ошибка сервера');

        const files = await response.json();

        if (!Array.isArray(files)) {
            throw new Error('Некорректный формат данных');
        }

        updateMap(map, files);
        updateFilesTable(files, map);
    } catch (error) {
        console.error('Ошибка:', error);
        updateFilesTable([], map);
    }
}

function updateMap(map, files) {
    clearMap(map);

    // Группируем файлы по регионам
    const filesByRegion = {};
    files.forEach(file => {
        if (!filesByRegion[file.Region]) {
            filesByRegion[file.Region] = [];
        }
        filesByRegion[file.Region].push(file);
    });

    // Добавляем маркеры для каждого региона
    for (const region in filesByRegion) {
        const regionFiles = filesByRegion[region];
        const regionGroup = L.layerGroup();

        regionFiles.forEach(file => {
            const marker = L.marker([file.Lat, file.Lon], {
                title: file.Filename
            }).bindPopup(`
                <b>${file.Filename}</b><br>
                <i>${file.Region}</i><br>
                Дата: ${new Date(file.Date).toLocaleDateString()}<br>
                <a href="${file.Path}" target="_blank">Открыть файл</a>
            `);

            regionGroup.addLayer(marker);
        });

        regionGroup.addTo(map);
        regionLayers[region] = regionGroup;
    }

    // Автоматически подбираем масштаб
    if (files.length > 0) {
        const bounds = L.latLngBounds(files.map(f => [f.Lat, f.Lon]));
        map.fitBounds(bounds, { padding: [50, 50] });
    }
}

function updateFilesTable(files, map) {
    const tbody = document.getElementById('data-body');
    if (!tbody) return;

    tbody.innerHTML = '';

    const safeFiles = Array.isArray(files) ? files : [];

    safeFiles.forEach(file => {
        const row = document.createElement('tr');

        // Название файла
        const nameCell = document.createElement('td');
        nameCell.textContent = file.Filename;
        row.appendChild(nameCell);

        // Регион
        const regionCell = document.createElement('td');
        regionCell.textContent = file.Region;
        row.appendChild(regionCell);

        // Дата
        const dateCell = document.createElement('td');
        dateCell.textContent = new Date(file.Date).toLocaleDateString();
        row.appendChild(dateCell);

        // Действия
        const actionCell = document.createElement('td');
        const link = document.createElement('a');
        link.href = file.Path;
        link.textContent = 'Открыть';
        link.target = '_blank';
        actionCell.appendChild(link);

        // Кнопка для показа на карте
        const showOnMapBtn = document.createElement('button');
        showOnMapBtn.textContent = 'Показать на карте';
        showOnMapBtn.addEventListener('click', () => {
            map.setView([file.Lat, file.Lon], 10);
            const marker = Object.values(regionLayers)
                .flatMap(layer => layer.getLayers())
                .find(m => m.options.title === file.Filename);
            if (marker) marker.openPopup();
        });
        actionCell.appendChild(showOnMapBtn);

        row.appendChild(actionCell);

        // Подсветка при наведении
        row.addEventListener('mouseenter', () => {
            const marker = Object.values(regionLayers)
                .flatMap(layer => layer.getLayers())
                .find(m => m.options.title === file.Filename);
            if (marker) {
                marker.setIcon(L.icon({
                    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon-2x.png',
                    iconSize: [25, 41],
                    iconAnchor: [12, 41],
                    popupAnchor: [1, -34],
                }));
            }
        });

        row.addEventListener('mouseleave', () => {
            const marker = Object.values(regionLayers)
                .flatMap(layer => layer.getLayers())
                .find(m => m.options.title === file.Filename);
            if (marker) {
                marker.setIcon(L.icon({
                    iconUrl: 'https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.7.1/images/marker-icon.png',
                    iconSize: [25, 41],
                    iconAnchor: [12, 41],
                    popupAnchor: [1, -34],
                }));
            }
        });

        tbody.appendChild(row);
    });
}