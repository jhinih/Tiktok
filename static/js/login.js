fetch('/api/post', {
    headers: {
        'Authorization': 'Bearer ' + localStorage.getItem('jwtToken')
    }
})