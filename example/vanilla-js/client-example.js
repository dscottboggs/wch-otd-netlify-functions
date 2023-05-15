
function applyToAll(selector, cb) {
    document.querySelectorAll('.wch-otd-' + selector).forEach(cb)
}
function recursivelyDisplay(element) {
    element.style.display = element.dataset.showAsDisplay ?? 'auto'
    for (var child of element.children) {
        recursivelyDisplay(child)
    }
}
document.addEventListener('readystatechange', async function () {
    if (document.readyState === "interactive") {
        console.log('started')
        const result = await fetch('https://stories.workingclasshistory.com/api/v1/one_random_from_today'),
            { ok } = result,
            data = await result.json()
        if (ok) {
            applyToAll('error', el => el.style.display = 'none')
            applyToAll('title', el => el.innerText += data.title)
            applyToAll('content', el => el.innerHTML = data.content)
            applyToAll('excerpt', el => el.innerHTML = data.excerpt)
            applyToAll('story-link', el => el.href = data.url)
            if (data.media) {
                applyToAll('media-placeholder', el => {
                    const img = document.createElement('img')
                    img.src = data.media.url
                    img.alt = data.media.caption
                    el.appendChild(img)
                })
            }
            console.log('done')
        } else {
            const { error } = data
            applyToAll('error', errDisplay => {
                errDisplay.innerText = error
            })
            applyToAll('loaded-wrapper', el => el.style.display = 'none')
            console.error('error', error)
        }
    }
})