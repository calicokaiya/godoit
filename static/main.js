$(document).ready(function () {
    $('#createTaskForm').submit(function (event) {
        event.preventDefault(); // Prevent the default form submission

        // Serialize the form data
        var formData = $(this).serialize();
        console.log(formData)

        // Send an AJAX request
        $.ajax({
            url: '/tasks/create',
            type: 'POST',
            data: formData,
            success: function (response) {
                // Handle the successful response
                console.log(response);
                location.reload();
            },
            error: function (xhr, status, error) {
                // Handle the error
                console.error(error);
            }
        });
    });
    // Handle the click event on the Update button
    $('.updateButton').click(function () {
        var taskId = $(this).data('task-id'); // Get the task ID from the button's data attribute

        // Populate the modal fields with existing data
        var $row = $('#' + taskId);
        var title = $row.find('td:eq(0)').text();
        var description = $row.find('td:eq(1)').text();
        var dueDate = $row.find('td:eq(2)').text();

        $('#recipient-name').val(title);
        $('#message-text').val(description);
        $('#dueDate').val(dueDate);
        $('#taskId').val(taskId);

        // Show the modal
        $('#updateModal').modal('show');
    });

    // Handle the form submission
    $('#updateTaskForm').submit(function (event) {
        event.preventDefault();

        var formData = $(this).serialize();
        //var taskId = $('#taskId').val();
        console.log(formData)

        // Send an AJAX request to update the task
        $.ajax({
            url: '/tasks/update', // Adjust the URL to match your backend endpoint
            type: 'POST',
            data: formData,
            success: function (response) {
                console.log(response);
                location.reload(); // Reload the page or update the table as needed
            },
            error: function (xhr, status, error) {
                console.error(error);
            }
        });

        // Hide the modal
        $('#updateModal').modal('hide');
    });
    $('.deleteButton').click(function () {
        var taskId = $(this).data('task-id');
        console.log(taskId);
        var formData = $(this).serialize();

        // Send an AJAX request
        $.ajax({
            url: '/tasks/delete',
            type: 'POST',
            data: "id=" + taskId,
            success: function (response) {
                // Handle the successful response
                console.log(response);
                location.reload();
            },
            error: function (xhr, status, error) {
                // Handle the error
                console.error(error);
            }
        });
    });
});