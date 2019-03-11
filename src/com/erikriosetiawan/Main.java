package com.erikriosetiawan;

import java.util.Scanner;
import java.util.Calendar;
import java.text.SimpleDateFormat;

public class Main {

    private static String nama;
    private static String kodeBuku;
    private static int jumlahBuku;
    private static Calendar myCalendar = Calendar.getInstance();
    private static SimpleDateFormat dateFormat = new SimpleDateFormat("dd MMMM yyyy HH:mm:ss");
    public static void main(String[] args) {

        Scanner inputUser = new Scanner(System.in);

        print("========== PROGRAM PEMINJAMAN BUKU PERPUSTAKAAN ==========");
        print("\n");
        print("Masukkan Nama Anda : ");
        nama = inputUser.nextLine();
        print("-----> Kode Buku <-----\n");
        print("Matematika : A1B1\nBahasa Indonesia : A2B3\nBahasa Inggris : A4B6\nIPA : A8B3\n");
        print("Masukkan Kode Buku : ");
        kodeBuku = inputUser.nextLine();
        print("Masukkan Jumlah Buku : ");
        jumlahBuku =inputUser.nextInt();

        output();
    }

    private static void print(String str) {
        System.out.print(str);
    }

    private static void output() {
        print("\n");
        print("========== DATA PEMINJAM ==========\n");
        print("Nama : " + nama + "\n");
        if (kodeBuku.equalsIgnoreCase("A1B1")) {
            print("Nama Buku : Matematika");
        } else if (kodeBuku.equalsIgnoreCase("A2B3")) {
            print("Nama Buku : Bahasa Indonesia");
        } else if (kodeBuku.equalsIgnoreCase("A4B6")) {
            print("Nama Buku : Bahasa Inggris");
        } else if (kodeBuku.equalsIgnoreCase("A8B3")) {
            print("Nama Buku : IPA");
        } else {
            print("Kode Buku Tidak Valid!");
        }
        print("\n");
        print("Jumlah Buku : " + jumlahBuku + "\n");
        myCalendar.add(Calendar.DAY_OF_MONTH, 7);
        print("Tanggal Pengembalian : " + dateFormat.format(myCalendar.getTime()));

    }


}

