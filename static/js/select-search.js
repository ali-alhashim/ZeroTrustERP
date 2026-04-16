document.addEventListener("DOMContentLoaded", function() {
    document.querySelectorAll('.select-search').forEach(select => {
        select.style.display = 'none';

        const container = document.createElement('div');
        container.className = 'select-search-container';
        
        const tagArea = document.createElement('div');
        tagArea.style.display = 'contents'; // Let tags mingle with input
        
        const input = document.createElement('input');
        input.className = 'select-search-input';
        input.placeholder = "Search...";

        const dropdown = document.createElement('div');
        dropdown.className = 'select-search-dropdown';

        container.append(tagArea, input, dropdown);
        select.parentNode.insertBefore(container, select);

        const renderTags = () => {
            tagArea.innerHTML = '';
            Array.from(select.selectedOptions).forEach(opt => {
                const tag = document.createElement('span');
                tag.className = 'select-tag';
                tag.innerHTML = `${opt.text} <span class="remove-tag" data-val="${opt.value}">&times;</span>`;
                
                tag.querySelector('.remove-tag').onclick = (e) => {
                    e.stopPropagation();
                    opt.selected = false;
                    renderTags();
                    select.dispatchEvent(new Event('change'));
                };
                tagArea.appendChild(tag);
            });
        };

        const updateDropdown = () => {
            dropdown.innerHTML = '';
            const filter = input.value.toLowerCase();
            let hasResults = false;

            Array.from(select.options).forEach(opt => {
                if (opt.text.toLowerCase().includes(filter)) {
                    const item = document.createElement('div');
                    item.className = 'select-search-option' + (opt.selected ? ' selected' : '');
                    item.style.padding = '8px';
                    item.style.cursor = 'pointer';
                    item.textContent = opt.text;
                    
                    item.onclick = () => {
                        opt.selected = !opt.selected;
                        renderTags();
                        input.value = '';
                        dropdown.style.display = 'none';
                        select.dispatchEvent(new Event('change'));
                    };
                    dropdown.appendChild(item);
                    hasResults = true;
                }
            });
            dropdown.style.display = hasResults ? 'block' : 'none';
        };

        input.onfocus = updateDropdown;
        input.oninput = updateDropdown;
        
        // Close if clicking away
        document.addEventListener('click', (e) => {
            if (!container.contains(e.target)) dropdown.style.display = 'none';
        });

        // Initialize tags if options are already selected (on edit page)
        renderTags();
    });
});