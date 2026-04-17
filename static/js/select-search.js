document.addEventListener("DOMContentLoaded", function() {
    document.querySelectorAll('.select-search').forEach(select => {
        select.style.display = 'none';

        const container = document.createElement('div');
        container.className = 'select-search-container';
        
        const tagArea = document.createElement('div');
        tagArea.className = 'tag-area';
        
        const input = document.createElement('input');
        input.className = 'select-search-input';
        input.placeholder = "Click to view or type...";

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

        const addCustomValue = () => {
            const val = input.value.trim().replace(/,+$/, '');
            if (val === '') return;

            let existingOpt = Array.from(select.options).find(
                opt => opt.text.toLowerCase() === val.toLowerCase()
            );

            if (existingOpt) {
                existingOpt.selected = true;
            } else {
                const newOpt = new Option(val, val, true, true);
                select.add(newOpt);
            }

            input.value = '';
            renderTags();
            dropdown.style.display = 'none';
            select.dispatchEvent(new Event('change'));
        };

        const updateDropdown = () => {
            dropdown.innerHTML = '';
            const filter = input.value.toLowerCase().trim();
            let hasResults = false;

            Array.from(select.options).forEach(opt => {
                // Show if filter matches OR if filter is empty (show all)
                if (!filter || opt.text.toLowerCase().includes(filter)) {
                    const item = document.createElement('div');
                    item.className = 'select-search-option' + (opt.selected ? ' selected' : '');
                    item.textContent = opt.text;
                    
                    item.onclick = (e) => {
                        e.stopPropagation();
                        opt.selected = !opt.selected;
                        input.value = '';
                        renderTags();
                        // Keep dropdown open for multi-select, or hide it:
                        // dropdown.style.display = 'none'; 
                        updateDropdown(); // Refresh highlights
                        select.dispatchEvent(new Event('change'));
                    };
                    dropdown.appendChild(item);
                    hasResults = true;
                }
            });
            dropdown.style.display = hasResults ? 'block' : 'none';
        };

        input.addEventListener('blur', () => {
            setTimeout(addCustomValue, 200);
        });

        input.addEventListener('input', updateDropdown);
        input.onfocus = updateDropdown; // Shows all options on click
        
        document.addEventListener('click', (e) => {
            if (!container.contains(e.target)) dropdown.style.display = 'none';
        });

        renderTags();
    });
});