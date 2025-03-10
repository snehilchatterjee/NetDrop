package com.primex.netdrop;

import android.os.Bundle;
import android.os.Handler;
import android.os.Looper;
import androidx.appcompat.app.AppCompatActivity;
import com.primex.netdrop.databinding.ActivityMainBinding;
import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.PrintWriter;
import java.net.Socket;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class MainActivity extends AppCompatActivity {
    private ActivityMainBinding binding;
    private final ExecutorService executorService = Executors.newSingleThreadExecutor();
    private final Handler mainHandler = new Handler(Looper.getMainLooper());

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        binding = ActivityMainBinding.inflate(getLayoutInflater());
        setContentView(binding.getRoot());

        binding.btnSend.setOnClickListener(v -> sendMessage(binding.edtSendMessage.getText().toString()));
    }

    private void sendMessage(final String msg) {
        executorService.execute(() -> {
            try (Socket socket = new Socket("10.0.2.15", 5000);
                 OutputStream out = socket.getOutputStream();
                 PrintWriter output = new PrintWriter(out, true);
                 BufferedReader input = new BufferedReader(new InputStreamReader(socket.getInputStream()))) {

                output.println(msg);
                final String response = input.readLine();

                mainHandler.post(() -> {
                    if (response != null && !response.trim().isEmpty()) {
                        binding.tvReplyFromServer.append("\nFrom Server: " + response);
                    }
                });
            } catch (IOException e) {
                e.printStackTrace();
            }
        });
    }

    @Override
    protected void onDestroy() {
        super.onDestroy();
        executorService.shutdown();
    }
}