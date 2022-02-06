getConversations();

async function getConversations(){
    const response = await fetch("http://localhost:9191/api/v1/groups/conversations")
    data = await response.json();

    const tableData = data
        .map(function(value){
            return `<tr>
                <td>${value.conversationID}</td>
                <td>${value.posterName}</td>
                <td>${value.postContent}</td>
                <td>${value.postDate}</td>
                <td>${value.postTime}</td>
            </tr>`;
        }).join("");

    const tableBody = document.querySelector("#tableBody");
    tableBody.innerHTML = tableData;
}