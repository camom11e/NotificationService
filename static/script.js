document.getElementById('send-button').addEventListener('click', async () => {
 const data = {
        sms: document.querySelector('body > div.container > div.block.checkboxes > label:nth-child(2) > input[type=checkbox]').checked,
		email: document.querySelector('body > div.container > div.block.checkboxes > label:nth-child(3) > input[type=checkbox]').checked,
		inApp: document.querySelector('body > div.container > div.block.checkboxes > label:nth-child(4) > input[type=checkbox]').checked,
		push: document.querySelector('body > div.container > div.block.checkboxes > label:nth-child(5) > input[type=checkbox]').checked,
        message: document.querySelector('.message-settings textarea').value,
        dateTime: document.querySelector('.message-settings input[type="datetime-local"]').value,
    };
    alert("Отправляемые данные:"+ JSON.stringify(data, null, 2));

    try {
        const response = await fetch('/api/send', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data),
        });

        if (!response.ok) throw new Error('Ошибка сервера');
        alert('Данные отправлены!');
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Не удалось отправить данные');
    }
});